package clienthttp

import (
	"bytes"
	"io"
	"net/http"
)

type HTTPRequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	query   map[string]string
	body    io.Reader
}

func NewRequest(method, url string) HTTPRequestBuilder {
	return HTTPRequestBuilder{
		method:  method,
		url:     url,
		headers: make(map[string]string),
		query:   make(map[string]string),
	}
}

func (h HTTPRequestBuilder) WithHeader(key, value string) HTTPRequestBuilder {
	h.headers[key] = value
	return h
}

func (h HTTPRequestBuilder) WithQueryParam(key, value string) HTTPRequestBuilder {
	h.query[key] = value
	return h
}

func (h HTTPRequestBuilder) WithBodyBytes(body []byte) HTTPRequestBuilder {
	h.body = bytes.NewBuffer(body)
	return h
}

func (h HTTPRequestBuilder) Build() *http.Request {
	req, _ := http.NewRequest(h.method, h.url, h.body)

	for key, value := range h.headers {
		req.Header.Set(key, value)
	}

	q := req.URL.Query()
	for k, v := range h.query {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	return req
}
