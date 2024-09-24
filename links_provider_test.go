package main

import (
	"sync"
	"testing"
)

// TestURLGeneratorWithHardcodedURLs verifies that the URLGenerator outputs the exact expected URLs.
func TestURLGeneratorWithHardcodedURLs(t *testing.T) {
	numURLs := 10
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
		"http://example.com/entity-9-19.json",
	}
	urlChan := make(chan string, numURLs)
	var wg sync.WaitGroup

	wg.Add(1)
	go URLGenerator(urlChan, numURLs, &wg)

	// Wait for URLGenerator to finish
	wg.Wait()

	generatedURLs := make([]string, 0, numURLs)
	for url := range urlChan {
		generatedURLs = append(generatedURLs, url)
	}

	// Check if the generated URLs match the expected URLs
	if len(generatedURLs) != len(expectedURLs) {
		t.Errorf("Generated URLs count mismatch: expected %d, got %d", len(expectedURLs), len(generatedURLs))
	}

	for i, url := range generatedURLs {
		if url != expectedURLs[i] {
			t.Errorf("URL mismatch at index %d: expected %s, got %s", i, expectedURLs[i], url)
		}
	}

	// assert channel is closed closed
	_, ok := <-urlChan
	if ok {
		t.Errorf("Channel should be closed but was still open")
	}
}
