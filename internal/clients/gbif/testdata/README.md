# Animal occurrence fixtures

These are complete GBIF API response bodies captured on 2026-09-23 and formatted
as JSON. Each request uses the GBIF Backbone checklist and `limit=2` to keep the
fixtures small. Tests replay the responses locally; they do not call GBIF.

| Fixture | Request |
| --- | --- |
| `occurrences_humpback_whale.json` | [Humpback whale](https://api.gbif.org/v1/occurrence/search?checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=5220086&limit=2) |
| `occurrences_bottlenose_dolphin.json` | [Common bottlenose dolphin](https://api.gbif.org/v1/occurrence/search?checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=2440447&limit=2) |
| `occurrences_alligators.json` | [American and Chinese alligators](https://api.gbif.org/v1/occurrence/search?checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=2441370&taxonKey=2441368&limit=2) |
| `occurrences_marine_animals.json` | [Blue whale, humpback whale and common bottlenose dolphin](https://api.gbif.org/v1/occurrence/search?checklistKey=d7dddbf4-2cf0-4f39-9b2a-bb099caae36c&taxonKey=2440735&taxonKey=5220086&taxonKey=2440447&limit=2) |

The humpback response includes a subspecies (`taxonKey=7388406`) whose
`speciesKey=5220086` matches the requested species. Tests therefore check species
membership rather than assuming every returned taxon key equals a query key.
The combined queries need not return every selected species on the first page.

To refresh a fixture, download its linked request and format the response with
`python3 -m json.tool`. Review the returned records and update the corresponding
expectations in `occurrences_test.go` if necessary. Normal application startup
does not need to fetch or regenerate these files.
