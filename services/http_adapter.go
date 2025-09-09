package services

import (
	"mime/multipart"
	"net/http"
)

// HTTPRequestAdapter adapts *http.Request to our HTTPRequest interface
type HTTPRequestAdapter struct {
	req *http.Request
}

// NewHTTPRequestAdapter creates a new HTTP request adapter
func NewHTTPRequestAdapter(req *http.Request) HTTPRequest {
	return &HTTPRequestAdapter{req: req}
}

// ParseMultipartForm parses multipart form
func (h *HTTPRequestAdapter) ParseMultipartForm(maxMemory int64) error {
	return h.req.ParseMultipartForm(maxMemory)
}

// FormValue gets form value
func (h *HTTPRequestAdapter) FormValue(name string) string {
	return h.req.FormValue(name)
}

// FormFile gets form file
func (h *HTTPRequestAdapter) FormFile(name string) (File, FileHeader, error) {
	file, header, err := h.req.FormFile(name)
	if err != nil {
		return nil, nil, err
	}
	return &FileAdapter{file: file}, &FileHeaderAdapter{header: header}, nil
}

// FileAdapter adapts multipart.File to our File interface
type FileAdapter struct {
	file multipart.File
}

// Close closes the file
func (f *FileAdapter) Close() error {
	return f.file.Close()
}

// Read reads from the file
func (f *FileAdapter) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

// FileHeaderAdapter adapts *multipart.FileHeader to our FileHeader interface
type FileHeaderAdapter struct {
	header *multipart.FileHeader
}

// Filename returns the filename
func (f *FileHeaderAdapter) Filename() string {
	return f.header.Filename
}
