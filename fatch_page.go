package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

// fetchPageContent attempts to fetch content from a URL, with retries.
func fetchPageContent(appCtx *AppContext, url string) ([]byte, error) {
	var (
		err        error
		resp       *http.Response
		bodyBytes  []byte
		contentErr error
	)

	// Wait for a token from the limiter
	if appCtx.Limiter != nil {
		// Create a context with a timeout for the limiter
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = appCtx.Limiter.Wait(ctx); err != nil {
			return nil, errors.New("rate limiter error")
		}
	}

	for attempt := 0; attempt < appCtx.MaxRetries; attempt++ {
		// Create a context with a timeout for the HTTP request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create a new HTTP request with the context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.New("failed to create HTTP request")
		}

		resp, err = appCtx.Client.Do(req)
		if err != nil {
			// Network error, retry
			contentErr = errors.New("network error")
		} else {
			defer resp.Body.Close()

			// Check for non-200 status codes
			if resp.StatusCode != http.StatusOK {
				contentErr = errors.New("received non-200 status code")
			} else {
				// Read and return the response body
				bodyBytes, err = io.ReadAll(resp.Body)
				if err != nil {
					contentErr = errors.New("error reading response body")
				} else {
					// Successful response
					return bodyBytes, nil
				}
			}
		}

		// Apply exponential backoff before retrying
		if attempt+1 < appCtx.MaxRetries {
			time.Sleep(time.Duration(1<<attempt) * time.Second)
			continue
		} else {
			return nil, contentErr
		}
	}

	return nil, errors.New("max retries exceeded")
}
