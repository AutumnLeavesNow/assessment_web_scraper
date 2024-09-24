package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type AppContext struct {
	Client     *http.Client
	Limiter    *rate.Limiter
	MaxRetries int
	Config     *Config
	Failed     *DLQueue
}

func fetcherWorker(appCtx *AppContext, urlChan <-chan string, dataChan chan<- FetchedData, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urlChan {
		log.Printf("[Fetcher] Fetching URL: %s", url)
		bodyBytes, err := fetchPageContent(appCtx, url)
		if err != nil {
			// Handle error
			log.Printf("[Fetcher] Error fetching URL %s: %v", url, err)
			appCtx.Failed.Add(url)
			continue // Continue to the next URL
		}
		// Send fetched data to the data channel
		dataChan <- FetchedData{URL: url, BodyBytes: bodyBytes}
	}
}

func processorWorker(appCtx *AppContext, dataChan <-chan FetchedData, wg *sync.WaitGroup) {
	defer wg.Done()
	for fetchedData := range dataChan {
		processFetchedData(appCtx, fetchedData)
	}
}

func RunApp(appCtx *AppContext) {
	// TODO create dynamic worker pools
	// based on load, response time, etc

	urlChan := make(chan string, appCtx.Config.NumLinks)
	dataChan := make(chan FetchedData, appCtx.Config.NumLinks)
	var wg sync.WaitGroup

	wg.Add(1)
	go URLGenerator(urlChan, appCtx.Config.NumLinks, &wg)

	fetcherWG := &sync.WaitGroup{}
	fetcherWG.Add(appCtx.Config.NumFetchWorkers)
	for i := 0; i < appCtx.Config.NumFetchWorkers; i++ {
		go fetcherWorker(appCtx, urlChan, dataChan, fetcherWG)
	}

	wg.Add(appCtx.Config.NumProcessWorkers)
	for i := 0; i < appCtx.Config.NumProcessWorkers; i++ {
		go processorWorker(appCtx, dataChan, &wg)
	}

	fetcherWG.Wait()
	close(dataChan)

	wg.Wait()
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	config, err := ReadConfig("config.json")
	if err != nil {
		log.Printf("Error reading config: %v", err)
		os.Exit(1)
	}
	appCtx := &AppContext{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Limiter:    rate.NewLimiter(rate.Limit(config.RateLimit), config.BurstLimit),
		MaxRetries: config.MaxRetries,
		Config:     config,
		Failed:     NewDLQueue(),
	}
	RunApp(appCtx)

	log.Printf("Finished work")
}
