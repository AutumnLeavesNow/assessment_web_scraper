package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func flattenMapKeys(m map[string]struct{}, separator string) string {
	if len(m) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	return strings.Join(keys, separator)
}

func setupMockResponses() map[string]*http.Response {
	responses := make(map[string]*http.Response)

	expectedURLs := []string{
		"http://example.com/product-0.html",
		"http://example.com/entity-1-11.json",
		"http://example.com/product-2.html",
		"http://example.com/entity-3-13.json",
		"http://example.com/product-4.html",
		"http://example.com/entity-5-15.json",
		"http://example.com/product-6.html",
		"http://example.com/entity-7-17.json",
		"http://example.com/product-8.html",
		// "http://example.com/entity-9-19.json", leave for error
	}

	for _, url := range expectedURLs {
		var body string
		var contentType string

		switch url {
		case "http://example.com/product-0.html":
			body = `<div data-id="product-0-id">Product 0 Content</div>`
			contentType = "text/html"
		case "http://example.com/entity-1-11.json":
			body = `{"title": "Entity 1 Title"}`
			contentType = "application/json"
		case "http://example.com/product-2.html":
			body = `<div data-id="product-2-id">Product 2 Content</div>`
			contentType = "text/html"
		case "http://example.com/entity-3-13.json":
			body = `{"title": "Entity 3 Title"}`
			contentType = "application/json"
		case "http://example.com/product-4.html":
			body = `<div data-id="product-4-id">Product 4 Content</div>`
			contentType = "text/html"
		case "http://example.com/entity-5-15.json":
			body = `{"title": "Entity 5 Title"}`
			contentType = "application/json"
		case "http://example.com/product-6.html":
			body = `<div data-id="product-6-id">Product 6 Content</div>`
			contentType = "text/html"
		case "http://example.com/entity-7-17.json":
			body = `{"title": "Entity 7 Title"}`
			contentType = "application/json"
		case "http://example.com/product-8.html":
			body = `<div data-id="product-8-id">Product 8 Content</div>`
			contentType = "text/html"

		default:
			body = `Not Found`
			contentType = "text/plain"
		}

		responses[url] = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
		responses[url].Header.Set("Content-Type", contentType)
	}

	responses["http://example.com/entity-9-19.json"] = &http.Response{
		StatusCode: http.StatusForbidden,
	}

	return responses
}

type MockRoundTripper struct {
	responses map[string]*http.Response
}

func (mrt *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if resp, ok := mrt.responses[req.URL.String()]; ok {
		return resp, nil
	}
	return nil, fmt.Errorf("no mock response for URL: %s", req.URL.String())
}

func TestRunApp(t *testing.T) {
	// Setup mock responses
	mockResponses := setupMockResponses()

	// Create a mock HTTP client with the custom responses
	mockClient := &http.Client{
		Transport: &MockRoundTripper{responses: mockResponses},
		Timeout:   10 * time.Second,
	}

	config := &Config{
		NumLinks:          10,
		NumFetchWorkers:   5,
		NumProcessWorkers: 3,
		MaxRetries:        3,
		RateLimit:         1.0,
		BurstLimit:        5,
		Categories: []Category{
			{
				ID:   "entityDetails",
				Name: "entity",
				Path: "/entity-{slug}-{uuid}.json",
				Type: "json",
				DataFields: []DataField{
					{
						FieldName: "title",
						JSONPath:  "title",
					},
				},
			},
			{
				ID:   "productPage",
				Name: "product",
				Path: "/product-{slug}.html",
				Type: "html",
				DataFields: []DataField{
					{
						FieldName:   "dataId",
						CSSSelector: "div[data-id]",
						Attribute:   "data-id",
					},
				},
			},
		},
	}

	appCtx := &AppContext{
		Client:     mockClient,
		Limiter:    rate.NewLimiter(rate.Limit(config.RateLimit), config.BurstLimit),
		MaxRetries: config.MaxRetries,
		Config:     config,
		Failed:     NewDLQueue(),
	}

	// Run the application with the mocked context
	RunApp(appCtx)

	// assert failed urls
	if len(appCtx.Failed.m) != 1 || func() bool {
		_, exists := appCtx.Failed.m["http://example.com/entity-9-19.json"]
		return !exists
	}() {
		failed_urls := flattenMapKeys(appCtx.Failed.m, ",")
		t.Errorf("Unexpected failed urls: %s", failed_urls)
	}
}
