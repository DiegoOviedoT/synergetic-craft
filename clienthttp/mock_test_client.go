package clienthttp

import (
	"context"
	"net/http"
)

type Client struct {
	api ClientHTTP
}

func NewTestClient(api ClientHTTP) *Client {
	return &Client{
		api: api,
	}
}

func (t *Client) SendRequest() {
	ctx := context.Background()

	req := NewRequest(http.MethodGet, "/test/unit-test").
		WithHeader("test-unit", "test one").
		WithQueryParam("query-param", "param").
		WithBodyBytes([]byte("body byte")).
		Build()

	t.api.Do(ctx, req)
}
