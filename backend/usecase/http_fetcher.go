package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NewHTTPFetcher creates a Fetcher backed by the provided HTTP client.
// When client is nil, a default client with a sensible timeout is used.
func NewHTTPFetcher(client *http.Client) Fetcher {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &httpFetcher{client: client}
}

type httpFetcher struct {
	client *http.Client
}

func (f *httpFetcher) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	if f == nil || f.client == nil {
		return nil, fmt.Errorf("http fetcher is not configured")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for %s: %w", url, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request for %s: %w", url, err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	return resp.Body, nil
}
