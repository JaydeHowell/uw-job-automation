package jobrequest

import (
	"net/http"
	"time"
)

type WorkdayClient struct {
	httpClient  *http.Client
	baseURL     string
	contentType string
}

func NewWorkdayClient() WorkdayClient {
	return WorkdayClient{
		httpClient:  newHTTPClient(),
		baseURL:     "https://uw.wd5.myworkdayjobs.com/wday/cxs/uw/UWHires",
		contentType: "application/json",
	}
}

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 10
	transport.IdleConnTimeout = 90 * time.Second

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
	return client
}
