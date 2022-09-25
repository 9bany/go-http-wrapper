package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

type RequestProtocol interface{}

type request[T any] struct {
	target    Target
	transport http.RoundTripper
	debug     bool
}

func NewRequest[DataResponseType any](target Target) RequestProtocol {
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: 5000,
		}).DialContext,
	}
	return &request[DataResponseType]{target: target, transport: transport}
}

func (requestor *request[DataResponseType]) Execute(ctx context.Context) (*Response[DataResponseType], error) {
	httpResp, err := requestor.request(ctx)
	if err != nil {
		return nil, err
	}
	return requestor.responseProcessing(httpResp)
}

func (requestor *request[DataResponseType]) responseProcessing(resp *http.Response) (*Response[DataResponseType], error) {
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		return nil, errors.New("request failed")
	}

	if resp == nil || resp.Body == nil {
		return nil, errors.New("response empty")
	}

	defer resp.Body.Close()

	response := &Response[DataResponseType]{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}

	err := json.NewDecoder(resp.Body).Decode(&response.Data)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (requestor *request[DataResponseType]) request(ctx context.Context) (*http.Response, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	req, err = http.NewRequest(requestor.target.GetMethod(),
		requestor.target.GetEndpoint(),
		bytes.NewReader(requestor.target.GetBody()),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if requestor.debug {
			dump(req, resp)
		}
	}()

	headers := requestor.target.GetHeader()
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err = requestor.transport.RoundTrip(req.WithContext(ctx))

	return resp, err
}

func dump(req *http.Request, res *http.Response) {
	reqDump, _ := httputil.DumpRequest(req, true)
	respDump, _ := httputil.DumpResponse(res, false)

	prettyPrintDump("Request Details", reqDump)
	prettyPrintDump("Response Details", respDump)
}

func prettyPrintDump(heading string, data []byte) {
	const separatorWidth = 60

	fmt.Printf("\n\n%s", strings.ToUpper(heading))
	fmt.Printf("\n%s\n\n", strings.Repeat("-", separatorWidth))
	fmt.Print(string(data))
}
