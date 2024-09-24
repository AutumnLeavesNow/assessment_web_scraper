package main

import (
	"os"
	"testing"
)

// TestReadConfig tests the configuration file loading.
func TestReadConfig(t *testing.T) {
	// Create a temporary configuration for testing
	testConfig := `{
        "categories": [
            {
                "id": "testCategory",
                "name": "test",
                "path": "/test-path",
                "type": "json",
                "dataFields": [
                    {
                        "fieldName": "testField",
                        "jsonPath": "testPath"
                    }
                ]
            }
        ]
    }`

	// Write the temporary config to a file
	tempFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	_, err = tempFile.WriteString(testConfig)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	_, err = ReadConfig(tempFile.Name())
	if err != nil {
		t.Errorf("ReadConfig() error = %v, wantErr %v", err, false)
	}
}
