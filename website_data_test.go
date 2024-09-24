package main

import (
	"testing"
)

// TestExtractDataWithCSS tests the extractDataWithCSS function.
func TestExtractDataWithCSS(t *testing.T) {
	tests := []struct {
		name      string
		htmlData  []byte
		selector  string
		attribute string
		want      string
		wantErr   bool
	}{
		{
			name:     "Extract text content",
			htmlData: []byte(`<div><p>Hello World</p></div>`),
			selector: "p",
			want:     "Hello World",
			wantErr:  false,
		},
		{
			name:      "Extract attribute value",
			htmlData:  []byte(`<img src="image.png" alt="An image">`),
			selector:  "img",
			attribute: "src",
			want:      "image.png",
			wantErr:   false,
		},
		{
			name:     "Non-matching selector",
			htmlData: []byte(`<div><p>No match here</p></div>`),
			selector: "h1",
			want:     "",
			wantErr:  false,
		},
		{
			name:      "Non-existent attribute",
			htmlData:  []byte(`<div><a href="link.html">Link</a></div>`),
			selector:  "a",
			attribute: "title",
			want:      "",
			wantErr:   false,
		},
		{
			name:     "Empty selector",
			htmlData: []byte(`<div><p>Text</p></div>`),
			selector: "",
			want:     "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDataWithCSS(tt.htmlData, tt.selector, tt.attribute)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractDataWithCSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractDataWithCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtractDataWithJSON tests the extractDataWithJSON function.
func TestExtractDataWithJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData []byte
		jsonPath string
		want     string
		wantErr  bool
	}{
		{
			name:     "Extract simple value",
			jsonData: []byte(`{"name": "John Doe", "age": 30}`),
			jsonPath: "name",
			want:     "John Doe",
			wantErr:  false,
		},
		{
			name:     "Extract nested value",
			jsonData: []byte(`{"user": {"id": 123, "email": "john@example.com"}}`),
			jsonPath: "user.email",
			want:     "john@example.com",
			wantErr:  false,
		},
		{
			name:     "Non-existent JSONPath",
			jsonData: []byte(`{"name": "John Doe", "age": 30}`),
			jsonPath: "address",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Extract from array",
			jsonData: []byte(`{"users": [{"id": 1}, {"id": 2}, {"id": 3}]}`),
			jsonPath: "users.1.id",
			want:     "2",
			wantErr:  false,
		},
		{
			name:     "Extract value from missing array index",
			jsonData: []byte(`{"numbers": [10, 20, 30]}`),
			jsonPath: "numbers.5",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Empty JSONPath",
			jsonData: []byte(`{"key": "value"}`),
			jsonPath: "",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDataWithJSON(tt.jsonData, tt.jsonPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractDataWithJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractDataWithJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
