package httputil_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"workshop/internal/common/httputil"

	"github.com/stretchr/testify/require"
)

type payload struct {
	Name string `json:"name"`
}

func TestGetDecodesAndSendsHeaders(t *testing.T) {
	var gotPath, gotQuery, gotAuth, gotAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
		gotAuth, gotAccept = r.Header.Get("Authorization"), r.Header.Get("Accept")
		_, err := w.Write([]byte(`{"name":"Passer domesticus"}`))
		require.NoError(t, err)
	}))
	defer srv.Close()

	c := httputil.New(httputil.Config{
		BaseURL: srv.URL + "/v1",
		Headers: http.Header{"Authorization": {"Bearer tok"}},
	})

	got, err := c.Get[payload](t.Context(), "species/search", url.Values{"q": {"Passer"}})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Name != "Passer domesticus" {
		t.Errorf("Name = %q, want %q", got.Name, "Passer domesticus")
	}
	if want := "/v1/species/search"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	if want := "q=Passer"; gotQuery != want {
		t.Errorf("query = %q, want %q", gotQuery, want)
	}
	if want := "Bearer tok"; gotAuth != want {
		t.Errorf("Authorization = %q, want %q", gotAuth, want)
	}
	if want := "application/json"; gotAccept != want {
		t.Errorf("Accept = %q, want %q", gotAccept, want)
	}
}

func TestPostEncodesBody(t *testing.T) {
	var gotBody, gotContentType, gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			http.Error(w, "reading request body", http.StatusInternalServerError)
			return
		}
		gotBody, gotContentType, gotMethod = string(b), r.Header.Get("Content-Type"), r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := httputil.New(httputil.Config{BaseURL: srv.URL})

	// struct{} for T: no response body expected.
	if _, err := c.Post[struct{}](t.Context(), "occurrences", payload{Name: "x"}); err != nil {
		t.Fatalf("Post: %v", err)
	}

	if want := `{"name":"x"}`; gotBody != want {
		t.Errorf("body = %q, want %q", gotBody, want)
	}
	if want := "application/json"; gotContentType != want {
		t.Errorf("Content-Type = %q, want %q", gotContentType, want)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
}

func TestNon2xxReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no such species", http.StatusNotFound)
	}))
	defer srv.Close()

	c := httputil.New(httputil.Config{BaseURL: srv.URL})

	_, err := c.Get[payload](t.Context(), "species/0", nil)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	apiErr, ok := errors.AsType[*httputil.Error](err)
	if !ok {
		t.Fatalf("error is %T, want *httputil.Error", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Body != "no such species" {
		t.Errorf("Body = %q, want %q", apiErr.Body, "no such species")
	}
}

func TestPerRequestHeaderOverridesDefault(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-Tenant")
		_, err := w.Write([]byte(`{}`))
		require.NoError(t, err)
	}))
	defer srv.Close()

	c := httputil.New(httputil.Config{
		BaseURL: srv.URL,
		Headers: http.Header{"X-Tenant": {"default"}},
	})

	_, err := c.Do[payload](t.Context(), httputil.Request{
		Path:    "thing",
		Headers: http.Header{"X-Tenant": {"override"}},
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if got != "override" {
		t.Errorf("X-Tenant = %q, want %q", got, "override")
	}
}

func TestWrappedErrorsKeepTheirChain(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`not json`))
		require.NoError(t, err)
	}))
	defer srv.Close()

	c := httputil.New(httputil.Config{BaseURL: srv.URL})

	_, err := c.Get[payload](t.Context(), "thing", nil)
	if err == nil {
		t.Fatal("expected a decode error, got nil")
	}

	// The wrap helper adds context...
	if !strings.HasPrefix(err.Error(), "httputil: decoding GET response: ") {
		t.Errorf("error = %q, want the httputil decode prefix", err)
	}

	// ...while leaving the underlying error reachable.
	if _, ok := errors.AsType[*json.SyntaxError](err); !ok {
		t.Errorf("error %q does not unwrap to *json.SyntaxError", err)
	}
}

func TestContextCancellationPropagates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	c := httputil.New(httputil.Config{BaseURL: srv.URL})

	if _, err := c.Get[payload](ctx, "slow", nil); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}
