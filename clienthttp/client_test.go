package clienthttp_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"synergetic-craft/clienthttp"
	"testing"
	"time"
)

func TestClientHttp_Do(t *testing.T) {
	t.Run("should return ok when server resolve ok the request [SUCCESS]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			fmt.Fprintf(w, "ok")
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(http.DefaultClient, ts.URL)

		_, statusCode, err := client.Do(context.TODO(), req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("should return timeout when the request time exceeded [TIMEOUT]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			time.Sleep(200 * time.Millisecond)

			fmt.Fprintf(w, "ok")
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(&http.Client{Timeout: 100}, ts.URL)

		_, _, err := client.Do(context.TODO(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})
}

func TestClientHttp_DoWithTimeout(t *testing.T) {
	t.Run("should return ok when server resolve [SUCCESS]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			body, _ := io.ReadAll(r.Body)

			assert.Equal(t, "body byte", string(body))

			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, `{"status":"ok"}`)
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(http.DefaultClient, ts.URL)

		var out struct {
			Status string `json:"status"`
		}

		err := client.DoWithTimeout(context.TODO(), req, 100, http.StatusOK, &out)

		assert.NoError(t, err)
	})

	t.Run("should return error when server resolve [TIMEOUT]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			body, _ := io.ReadAll(r.Body)

			assert.Equal(t, "body byte", string(body))

			time.Sleep(200 * time.Millisecond)

			fmt.Fprintf(w, `{"status":"ok"}`)
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(&http.Client{Timeout: 1000 * time.Second}, ts.URL)

		var out struct {
			Status string `json:"status"`
		}

		err := client.DoWithTimeout(context.TODO(), req, 100, http.StatusOK, &out)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("should return error when server resolve [UNMARSHALL]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			body, _ := io.ReadAll(r.Body)

			assert.Equal(t, "body byte", string(body))

			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, `{"status":error}`)
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(http.DefaultClient, ts.URL)

		var out struct {
			Status string `json:"status"`
		}

		err := client.DoWithTimeout(context.TODO(), req, 100, http.StatusOK, &out)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("should return error when server resolve with bad request [BAD REQUEST]", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/unit-test/test1", r.URL.Path)
			assert.Equal(t, "test one", r.Header.Get("test-unit"))
			assert.Equal(t, "query-param=param", r.URL.RawQuery)

			body, _ := io.ReadAll(r.Body)

			assert.Equal(t, "body byte", string(body))

			w.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(w, `{"status":"ok"}`)
		}))
		defer ts.Close()

		req := clienthttp.NewRequest(http.MethodGet, "/unit-test/test1").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte("body byte")).
			Build()

		client := clienthttp.NewClientHTTP(http.DefaultClient, ts.URL)

		var out struct {
			Status string `json:"status"`
		}

		err := client.DoWithTimeout(context.TODO(), req, 100, http.StatusCreated, &out)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status Code [ 400 ]")
	})
}
