package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	NumLinks          int        `json:"numLinks"`
	NumFetchWorkers   int        `json:"numFetchWorkers"`
	NumProcessWorkers int        `json:"numProcessWorkers"`
	MaxRetries        int        `json:"maxRetries"`
	RateLimit         float64    `json:"rateLimit"`
	BurstLimit        int        `json:"burstLimit"`
	Categories        []Category `json:"categories"`
}

type Category struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Type       string      `json:"type"`
	DataFields []DataField `json:"dataFields"`
}

type DataField struct {
	FieldName   string `json:"fieldName"`
	CSSSelector string `json:"cssSelector,omitempty"`
	JSONPath    string `json:"jsonPath,omitempty"`
	Attribute   string `json:"attribute,omitempty"`
}

func ReadConfig(configFile string) (*Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
