package services

import (
	"context"
)

// Database interface for database operations
type Database interface {
	WithContext(ctx context.Context) Transaction
}

// Transaction interface for database transactions
type Transaction interface {
	Begin() Tx
	Where(query interface{}, args ...interface{}) Query
	First(dest interface{}, conds ...interface{}) Query
}

// Tx interface for database transactions
type Tx interface {
	Create(value interface{}) Query
	CreateInBatches(value interface{}, batchSize int) Query
	Commit() Query
	Rollback()
	Error() error
}

// Query interface for database queries
type Query interface {
	Error() error
	Order(value interface{}) Query
	Find(dest interface{}) Query
	Count(count *int64) Query
	Distinct(column string) Query
	Pluck(column string, dest interface{}) Query
	First(dest interface{}, conds ...interface{}) Query
}

// HTTPRequest interface for HTTP request operations
type HTTPRequest interface {
	ParseMultipartForm(maxMemory int64) error
	FormValue(name string) string
	FormFile(name string) (File, FileHeader, error)
}

// File interface for uploaded files
type File interface {
	Close() error
	Read(p []byte) (n int, err error)
}

// FileHeader interface for file headers
type FileHeader interface {
	Filename() string
}

// ImportRequest represents the data needed for import
type ImportRequest struct {
	Title          string
	IgnoredHeaders []string
	FileContent    []byte
	Filename       string
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	ImportJobID    uint
	RequestCount   int
	UniqueDomains  int
	Summary        string
}
