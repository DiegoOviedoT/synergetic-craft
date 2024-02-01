package clienthttp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockClientAPI struct {
	mock.Mock
}

func (m *MockClientAPI) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type MockClient struct {
	mock *MockClientAPI
	ClientHTTP
	testing *testing.T
}

func NewMockClient(t *testing.T) *MockClient {
	mockClientAPI := new(MockClientAPI)

	return &MockClient{
		mock:       mockClientAPI,
		ClientHTTP: NewClientHTTP(mockClientAPI, ""),
		testing:    t,
	}
}

func (m *MockClient) Do(ctx context.Context, req *http.Request) (response []byte, statusCode int, err error) {
	return m.ClientHTTP.Do(ctx, req)
}

func (m *MockClient) DoWithTimeout(ctx context.Context, req *http.Request, timeout int64, expectedCode int, out interface{}) error {
	return m.ClientHTTP.DoWithTimeout(ctx, req, timeout, expectedCode, out)
}

func (m *MockClient) ExpectedRequest(expectedReq *http.Request, response []byte, statusCode int, APIErr error) *mock.Call {
	var (
		expected, actual string
		isMatch          bool
	)

	return m.mock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		isMatch = checkMethod(expectedReq, req) && checkURL(expectedReq, req)
		if !isMatch {
			m.logger("Method/Url don´t match",
				fmt.Sprintf("%s %s", expectedReq.Method, expectedReq.URL),
				fmt.Sprintf("%s %s", req.Method, req.URL))
			return false
		}

		expected, actual, isMatch = checkHeaders(expectedReq, req)
		if !isMatch {
			m.logger("Headers don´t match", expected, actual)
			return false
		}

		expected, actual, isMatch = checkQueryParams(expectedReq, req)
		if !isMatch {
			m.logger("Query Parameters don´t match", expected, actual)
			return false
		}

		expected, actual, isMatch = m.checkBody(expectedReq, req)
		if !isMatch {
			m.logger("Body don´t match", expected, actual)
			return false
		}

		return true
	})).Return(&http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(response))},
		APIErr,
	)
}

func (m *MockClient) AssertExpectations() {
	m.mock.AssertExpectations(m.testing)
}

func checkMethod(expected, req *http.Request) bool {
	return expected.Method == req.Method
}

func checkURL(expected, req *http.Request) bool {
	return expected.URL.Path == req.URL.Path
}

func checkHeaders(expected, req *http.Request) (string, string, bool) {
	var (
		expectedHeader, actualHeader string
	)

	if len(expected.Header) > 0 && len(req.Header) < 1 {
		return "Expected request with headers, but received request without any.",
			"Received request without any header.",
			false
	}

	for key := range req.Header {
		expectedHeader = expected.Header.Get(key)
		actualHeader = req.Header.Get(key)

		if expectedHeader != actualHeader {
			return msg(key, expectedHeader), msg(key, actualHeader), false
		}

	}

	return "", "", true
}

func checkQueryParams(expected, req *http.Request) (string, string, bool) {
	var expectedParam, actualParam string

	if len(expected.URL.Query()) > 0 && len(req.URL.Query()) < 1 {
		return "Expected request with query parameters, but received request without any.",
			"Received request without any query parameters.",
			false
	}

	for key := range req.URL.Query() {
		expectedParam = expected.URL.Query().Get(key)
		actualParam = req.URL.Query().Get(key)

		if expectedParam != actualParam {
			return msg(key, expectedParam), msg(key, actualParam), false
		}
	}

	return "", "", true
}

func (m *MockClient) checkBody(expected, req *http.Request) (string, string, bool) {
	expectedBody := m.readRequestBody(expected)
	actualBody := m.readRequestBody(req)

	if len(expectedBody) < 1 && len(actualBody) < 1 {
		return string(expectedBody), string(actualBody), true
	}

	if len(expectedBody) > 0 && len(actualBody) < 1 {
		return string(expectedBody), string(actualBody), false
	}

	if len(expectedBody) < 1 && len(actualBody) > 0 {
		return string(expectedBody), string(actualBody), false
	}

	if !assert.ObjectsAreEqual(string(expectedBody), string(actualBody)) {
		return string(expectedBody), string(actualBody), false
	}

	return "", "", true
}

func (m *MockClient) readRequestBody(req *http.Request) []byte {
	if req.Body == nil {
		return []byte{}
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		m.testing.Log("Error reading the request body: ", err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return body
}

func (m *MockClient) logger(msg string, expected, actual string) {
	m.testing.Fatal(
		fmt.Sprintf("%s: \nExpected: %s\nActual: %s", msg, expected, actual),
	)
}

func msg(key, value string) string {
	return fmt.Sprintf("key: \"%s\" value: \"%s\"", key, value)
}
