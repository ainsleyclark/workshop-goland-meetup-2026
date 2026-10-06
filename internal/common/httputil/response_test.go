package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"workshop/internal/common/httputil"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value any
		want  string
	}{
		{"unescaped HTML and URL", map[string]string{"url": "https://example.com/?a=1&b=<bird>"}, `{"url":"https://example.com/?a=1&b=<bird>"}`},
		{"nil", nil, `null`},
		{"string newline", "bird\n", `"bird\n"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			httputil.WriteJSON(rec, http.StatusCreated, tt.value)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "application/json", rec.Result().Header.Get("Content-Type"))
			assert.Equal(t, tt.want, rec.Body.String())
		})
	}
}

func TestWriteJSONEncodingFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	httputil.WriteJSON(rec, http.StatusCreated, struct {
		Name        string
		Unsupported chan int
	}{Name: "must not be written", Unsupported: make(chan int)})
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "text/plain; charset=utf-8", rec.Result().Header.Get("Content-Type"))
	assert.Equal(t, "json: unsupported type: chan int\n", rec.Body.String())
}

func TestWriteRawJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	raw := " {\"html\":\"<script>alert('bird')</script>\",\"url\":\"?a=1&b=2\"}\n"
	httputil.WriteRawJSON(rec, raw)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Result().Header.Get("Content-Type"))
	assert.Equal(t, "nosniff", rec.Result().Header.Get("X-Content-Type-Options"))
	assert.Equal(t, raw, rec.Body.String())
}
