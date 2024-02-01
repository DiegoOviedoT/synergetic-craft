package clienthttp

import (
	"net/http"
	"testing"
)

func TestMockClient_ExpectedRequest(t *testing.T) {
	t.Run("should return success when request is correct", func(t *testing.T) {
		expectedReq := NewRequest(http.MethodGet, "/test/unit-test").
			WithHeader("test-unit", "test one").
			WithQueryParam("query-param", "param").
			WithBodyBytes([]byte(`body byte`)).
			Build()

		f := setupTestMockClientFixture(t)
		f.mockClient.ExpectedRequest(expectedReq, []byte(``), 200, nil)

		f.resource.SendRequest()
		f.mockClient.AssertExpectations()
	})
}

type testMockClientFixture struct {
	mockClient *MockClient
	resource   *Client
}

func setupTestMockClientFixture(t *testing.T) *testMockClientFixture {
	mockClient := NewMockClient(t)

	return &testMockClientFixture{
		mockClient: mockClient,
		resource:   NewTestClient(mockClient),
	}
}
