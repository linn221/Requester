package services

import (
	"fmt"
	"io"
	"strings"
)

// FormParser handles parsing of multipart form data
type FormParser struct{}

// NewFormParser creates a new FormParser
func NewFormParser() *FormParser {
	return &FormParser{}
}

// ParseImportForm parses the import form data from HTTP request
func (p *FormParser) ParseImportForm(r HTTPRequest) (*ImportRequest, error) {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %v", err)
	}

	// Get form values
	title := r.FormValue("title")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	ignoredHeadersText := r.FormValue("ignoredHeaders")
	ignoredHeaders := strings.Fields(strings.ReplaceAll(ignoredHeadersText, "\n", " "))

	// Get uploaded file
	file, header, err := r.FormFile("harfile")
	if err != nil {
		return nil, fmt.Errorf("failed to get uploaded file: %v", err)
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return &ImportRequest{
		Title:          title,
		IgnoredHeaders: ignoredHeaders,
		FileContent:    fileContent,
		Filename:       header.Filename(),
	}, nil
}
