package clienthttp_test

import (
	"io"
	"net/http"
	"testing"

	"synergetic-craft/clienthttp"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRequestBuilder(t *testing.T) {
	t.Run("should return New Request [type=GET]", func(t *testing.T) {
		request := clienthttp.NewRequest(http.MethodGet, "/test/test-1").
			WithHeader("Test-one", "Test 1").
			WithQueryParam("limit", "10").
			WithQueryParam("offset", "0").
			WithBodyBytes([]byte(`Test unit`)).
			Build()

		body, _ := io.ReadAll(request.Body)

		assert.Equal(t, 1, len(request.Header))
		assert.Equal(t, "Test 1", request.Header.Get("Test-one"))
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, "limit=10&offset=0", request.URL.RawQuery)
		assert.Equal(t, "Test unit", string(body))
	})
}
