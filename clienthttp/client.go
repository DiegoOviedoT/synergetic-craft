package clienthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type ClientHTTP interface {
	Do(ctx context.Context, req *http.Request, timeout int64) (response []byte, statusCode int, err error)
	DoWithTimeout(ctx context.Context, req *http.Request, timeout int64, expectedCode int, out interface{}) error
}

type clientHttp struct {
	domain        string
	headerDefault http.Header
	client        *http.Client
}

func NewClientHTTP(domain string, timeout int64) *clientHttp {
	return &clientHttp{
		domain:        domain,
		headerDefault: make(http.Header),
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
				MaxConnsPerHost:     -1,
				MaxIdleConns:        100,
				DisableKeepAlives:   true,
			},
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *clientHttp) Do(ctx context.Context, req *http.Request) (response []byte, statusCode int, err error) {
	return c.do(ctx, req, 0)
}

func (c *clientHttp) DoWithTimeout(ctx context.Context, req *http.Request, timeout int64, expectedCode int, out interface{}) error {
	resp, statusCode, err := c.do(ctx, req, timeout)
	if err != nil {
		return err
	}

	if statusCode != expectedCode {
		return fmt.Errorf("status Code [ %d ], err: %v", statusCode, err)
	}

	err = json.Unmarshal(resp, &out)
	if err != nil {
		return err
	}

	return nil
}

func (c *clientHttp) do(ctx context.Context, req *http.Request, timeout int64) (response []byte, statusCode int, err error) {
	if timeout > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancel()

		ctx = ctxWithTimeout
	}

	url := c.domain + req.URL.Path
	if req.URL.RawQuery != "" {
		url += "?" + req.URL.RawQuery
	}

	req.URL, err = req.URL.Parse(url)
	if err != nil {
		return
	}

	for k := range c.headerDefault {
		req.Header.Add(k, c.headerDefault.Get(k))
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return body, resp.StatusCode, nil
}
