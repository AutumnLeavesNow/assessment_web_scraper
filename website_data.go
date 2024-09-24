package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

type FetchedData struct {
	URL       string
	BodyBytes []byte
}

func extractDataWithCSS(htmlData []byte, selector, attribute string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %w", err)
	}
	var data string
	if attribute == "" {
		data = doc.Find(selector).Text()
	} else {
		data, _ = doc.Find(selector).Attr(attribute)
	}
	return strings.TrimSpace(data), nil
}

func extractDataWithJSON(jsonData []byte, jsonPath string) (string, error) {
	jsonString := string(jsonData)
	result := gjson.Get(jsonString, jsonPath)
	if !result.Exists() {
		return "", fmt.Errorf("json path not found: %s", jsonPath)
	}
	return result.String(), nil
}

func processFetchedData(appCtx *AppContext, fetchedData FetchedData) {
	// TODO we should really be checking match for config path and not type
	log.Printf("[Processor] Processing URL: %s", fetchedData.URL)
	for _, category := range appCtx.Config.Categories {
		log.Printf("[Processor] Category: %s", category.Name)
		for _, field := range category.DataFields {
			var value string
			var err error
			if strings.HasSuffix(fetchedData.URL, ".html") && category.Type == "html" {
				value, err = extractDataWithCSS(fetchedData.BodyBytes, field.CSSSelector, field.Attribute)
			} else if strings.HasSuffix(fetchedData.URL, ".json") && category.Type == "json" {
				value, err = extractDataWithJSON(fetchedData.BodyBytes, field.JSONPath)
			} else {
				continue
			}
			if err != nil {
				log.Printf("[Processor] Error extracting data: %v", err)
				continue
			}
			log.Printf("[Processor]  %s: %v", field.FieldName, value)
			// TODO: Store the extracted data to a database
		}
	}
}
