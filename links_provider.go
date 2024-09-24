package main

import (
	"fmt"
	"log"
	"sync"
)

func URLGenerator(urlChan chan<- string, numURLs int, wg *sync.WaitGroup) {
	defer wg.Done()
	baseURLs := []string{
		"http://example.com/product-%d.html",
		"http://example.com/entity-%d-%d.json",
	}

	log.Println("Starting URL generation")
	for i := 0; i < numURLs; i++ {
		var url string
		if i%2 == 0 {
			url = fmt.Sprintf(baseURLs[0], i)
			urlChan <- url
			log.Printf("Generated URL: %s", url)
		} else {
			url = fmt.Sprintf(baseURLs[1], i, i+numURLs)
			urlChan <- url
			log.Printf("Generated URL: %s", url)
		}
	}

	log.Println("Finished generating URLs, closing channel")
	close(urlChan)
}
