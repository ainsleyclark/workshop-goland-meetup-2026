package gbif

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"workshop/internal/clients/gbif"
	"workshop/internal/domain/species"
)

// perSeed caps the occurrences fetched per seed search. Raise it for
// more data, keeping the file comfortably under GitHub's 100MB limit.
const perSeed = 1000

type scrapeConfig struct {
	BaseURL string // The GBIF API to copy.
	PerSeed int    // Maximum occurrences per seed search.
	Out     string // Where to write the gzipped SQL.
}

// dataPath finds the embedded data file from the starter, the mockapi
// folder or the repository root, whichever the scrape is run from.
func dataPath() string {
	const path = "internal/gbif/data/gbif.sql.gz"
	for _, dir := range []string{".", "cmd/mockapi", "starter/cmd/mockapi"} {
		p := filepath.Join(dir, path)
		if _, err := os.Stat(filepath.Dir(p)); err == nil {
			return p
		}
	}
	return path
}

// backbone is the GBIF Backbone checklist, which the taxon keys belong to.
const backbone = "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c"

// seed is one occurrence search whose results go into the mock.
type seed struct {
	name  string
	query url.Values
}

// seeds covers every animal the app knows, some broad groups so other
// searches find something, and the countries the themes are built around.
func seeds() []seed {
	var out []seed
	taxa := func(name string, keys ...int) {
		q := url.Values{"checklistKey": {backbone}}
		for _, k := range keys {
			q.Add("taxonKey", strconv.Itoa(k))
		}
		out = append(out, seed{name: name, query: q})
	}

	for _, a := range species.Animals() {
		taxa(a.String(), a.TaxonKeys()...)
	}
	taxa("mammals", 359)
	taxa("birds", 212)
	taxa("reptiles", 358)
	taxa("amphibians", 131)
	taxa("ray-finned fish", 204)
	taxa("insects", 216)
	taxa("plants", 6)
	taxa("fungi", 5)

	for _, country := range []string{"KE", "AQ", "GB", "US", "AU", "MX"} {
		out = append(out, seed{name: "country " + country, query: url.Values{"country": {country}}})
	}

	return out
}

// Scrape copies the live GBIF API into the data file the mock embeds.
func Scrape(ctx context.Context, logger *slog.Logger) error {
	return scrape(ctx, logger, scrapeConfig{
		BaseURL: gbif.BaseURL,
		PerSeed: perSeed,
		Out:     dataPath(),
	})
}

func scrape(ctx context.Context, logger *slog.Logger, cfg scrapeConfig) error {
	c := &upstream{baseURL: cfg.BaseURL, http: &http.Client{Timeout: time.Minute}}
	d := newDataset()

	for _, s := range seeds() {
		before := len(d.Occurrences)
		for offset := 0; offset < cfg.PerSeed; offset += maxLimit {
			limit := min(maxLimit, cfg.PerSeed-offset)
			q := url.Values{}
			for k, v := range s.query {
				q[k] = v
			}
			q.Set("limit", strconv.Itoa(limit))
			q.Set("offset", strconv.Itoa(offset))

			body, err := c.get(ctx, "/occurrence/search", q)
			if err != nil {
				return fmt.Errorf("scraping %s: %w", s.name, err)
			}
			n, err := d.addPage(body)
			if err != nil {
				return fmt.Errorf("scraping %s: %w", s.name, err)
			}
			if n < limit {
				break // No more records.
			}
		}
		logger.Info("Scraped occurrences", "seed", s.name, "new", len(d.Occurrences)-before, "total", len(d.Occurrences))
	}

	keys := d.speciesKeys()
	logger.Info("Scraping species and media", "species", len(keys))
	if err := c.species(ctx, logger, d, keys); err != nil {
		return err
	}

	// Writes are scoped to the data folder, so nothing outside it
	// can be touched whatever the path says.
	root, err := os.OpenRoot(filepath.Dir(cfg.Out))
	if err != nil {
		return err
	}
	defer root.Close() //nolint:errcheck

	name := filepath.Base(cfg.Out)
	f, err := root.Create(name + ".tmp")
	if err != nil {
		return err
	}
	if err = d.writeSQL(f); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = root.Rename(name+".tmp", name); err != nil {
		return err
	}

	info, err := root.Stat(name)
	if err != nil {
		return err
	}
	logger.Info("Wrote GBIF data",
		"file", cfg.Out,
		"size", fmt.Sprintf("%.1fMB", float64(info.Size())/(1<<20)),
		"occurrences", len(d.Occurrences),
		"species", len(d.Species),
	)
	return nil
}

// species fetches each species and its media, a few at a time,
// as there are far more of them than occurrence pages.
func (c *upstream) species(ctx context.Context, logger *slog.Logger, d *dataset, keys []int) error {
	const workers = 8

	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		firstErr error
		done     int
	)
	jobs := make(chan int)
	for range workers {
		wg.Go(func() {
			for key := range jobs {
				path := "/species/" + strconv.Itoa(key)
				sp, err := c.get(ctx, path, nil)
				var media []byte
				if err == nil {
					media, err = c.get(ctx, path+"/media", nil)
				}

				mu.Lock()
				switch {
				case errors.Is(err, errNotFound):
				case err != nil && firstErr == nil:
					firstErr = fmt.Errorf("species %d: %w", key, err)
				case err != nil:
				default:
					d.Species[key] = sp
					d.Media[key] = media
				}
				if done++; done%250 == 0 {
					logger.Info("Scraped species", "done", done, "of", len(keys))
				}
				mu.Unlock()
			}
		})
	}

	for _, key := range keys {
		mu.Lock()
		failed := firstErr != nil
		mu.Unlock()
		if failed || ctx.Err() != nil {
			break
		}
		jobs <- key
	}
	close(jobs)
	wg.Wait()

	return errors.Join(firstErr, ctx.Err())
}

// upstream fetches raw response bodies from the real API.
// It deliberately skips the typed gbif client, since decoding and
// re-encoding would lose fields and change how values are written.
type upstream struct {
	baseURL string
	http    *http.Client
}

var errNotFound = errors.New("not found")

// get retries rate limits and server errors with a growing wait.
func (c *upstream) get(ctx context.Context, path string, q url.Values) ([]byte, error) {
	uri := c.baseURL + path
	if len(q) > 0 {
		uri += "?" + q.Encode()
	}

	var lastErr error
	for attempt := range 5 {
		if attempt > 0 {
			select {
			case <-time.After(time.Second << (attempt - 1)):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		body, status, err := c.do(ctx, uri)
		switch {
		case err != nil:
			lastErr = err
		case status == http.StatusNotFound:
			return nil, errNotFound
		case status == http.StatusTooManyRequests || status >= 500:
			lastErr = fmt.Errorf("GET %s: %s", uri, http.StatusText(status))
		case status != http.StatusOK:
			return nil, fmt.Errorf("GET %s: %s", uri, http.StatusText(status))
		default:
			return body, nil
		}
	}
	return nil, lastErr
}

func (c *upstream) do(ctx context.Context, uri string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(res.Body)
	return body, res.StatusCode, err
}
