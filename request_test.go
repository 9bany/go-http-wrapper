package request

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	dataMocking string
	statusCode  int
}

func newMockTransport(data string, statusCode int) http.RoundTripper {
	return &mockTransport{
		dataMocking: data,
		statusCode:  statusCode,
	}
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: t.statusCode,
	}
	response.Header.Set("Content-Type", "application/json")

	responseBody := t.dataMocking
	response.Body = io.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

type exampleTarget struct{}

func (target *exampleTarget) GetMethod() string {
	return ""
}

func (target *exampleTarget) GetEndpoint() string {
	return ""
}

func (target *exampleTarget) GetBody() []byte {
	return nil
}

func (target *exampleTarget) GetHeader() Header {
	return map[string]string{}
}

func TestExecuteRequest(t *testing.T) {
	type Data struct {
		Msg string `json:"data"`
	}

	testCases := []struct {
		name          string
		dataMocking   string
		statusCode    int
		checkResponse func(*Response[Data], error)
	}{
		{
			name: "OK",
			dataMocking: `{
				"data": "data_message"
			}`,
			statusCode: http.StatusOK,
			checkResponse: func(resp *Response[Data], err error) {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "data_message", resp.Data.Msg)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "Status code error",
			dataMocking: `{
				"data": "data_message"
			}`,
			statusCode: http.StatusInternalServerError,
			checkResponse: func(resp *Response[Data], err error) {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "Empty body error",
			dataMocking: ``,
			statusCode: http.StatusOK,
			checkResponse: func(resp *Response[Data], err error) {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target := &exampleTarget{}
			requestor := &request[Data]{
				target:    target,
				transport: newMockTransport(tc.dataMocking, tc.statusCode),
			}
			resp, err := requestor.Execute(context.Background())
			tc.checkResponse(resp, err)
		})
	}
}
