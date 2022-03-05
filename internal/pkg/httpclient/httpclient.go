package httpclient

import (
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	timeout time.Duration
}

func New(timeout time.Duration) *HttpClient {
	return &HttpClient{
		timeout: timeout,
	}
}

func (h *HttpClient) Download(endpoint string) (io.ReadCloser, error) {
	httpClient := http.Client{
		Timeout: h.timeout,
	}

	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
