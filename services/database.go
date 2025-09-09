package services

import (
	"context"
	"gorm.io/gorm"
)

// GormDatabaseAdapter adapts GORM to our Database interface
type GormDatabaseAdapter struct {
	db *gorm.DB
}

// NewGormDatabaseAdapter creates a new GORM database adapter
func NewGormDatabaseAdapter(db *gorm.DB) Database {
	return &GormDatabaseAdapter{db: db}
}

// WithContext returns a new database instance with context
func (g *GormDatabaseAdapter) WithContext(ctx context.Context) Transaction {
	return &GormTransactionAdapter{db: g.db.WithContext(ctx)}
}

// GormTransactionAdapter adapts GORM transaction to our Transaction interface
type GormTransactionAdapter struct {
	db *gorm.DB
}

// Begin starts a new transaction
func (g *GormTransactionAdapter) Begin() Tx {
	return &GormTxAdapter{db: g.db.Begin()}
}

// Where adds a WHERE clause
func (g *GormTransactionAdapter) Where(query interface{}, args ...interface{}) Query {
	return &GormQueryAdapter{db: g.db.Where(query, args...)}
}

// First finds the first record
func (g *GormTransactionAdapter) First(dest interface{}, conds ...interface{}) Query {
	return &GormQueryAdapter{db: g.db.First(dest, conds...)}
}

// GormTxAdapter adapts GORM transaction to our Tx interface
type GormTxAdapter struct {
	db *gorm.DB
}

// Create creates a record
func (g *GormTxAdapter) Create(value interface{}) Query {
	return &GormQueryAdapter{db: g.db.Create(value)}
}

// CreateInBatches creates records in batches
func (g *GormTxAdapter) CreateInBatches(value interface{}, batchSize int) Query {
	return &GormQueryAdapter{db: g.db.CreateInBatches(value, batchSize)}
}

// Commit commits the transaction
func (g *GormTxAdapter) Commit() Query {
	return &GormQueryAdapter{db: g.db.Commit()}
}

// Rollback rolls back the transaction
func (g *GormTxAdapter) Rollback() {
	g.db.Rollback()
}

// Error returns the error from the transaction
func (g *GormTxAdapter) Error() error {
	return g.db.Error
}

// GormQueryAdapter adapts GORM query to our Query interface
type GormQueryAdapter struct {
	db *gorm.DB
}

// Error returns the error from the query
func (g *GormQueryAdapter) Error() error {
	return g.db.Error
}

// Order adds an ORDER BY clause
func (g *GormQueryAdapter) Order(value interface{}) Query {
	return &GormQueryAdapter{db: g.db.Order(value)}
}

// Find finds records
func (g *GormQueryAdapter) Find(dest interface{}) Query {
	return &GormQueryAdapter{db: g.db.Find(dest)}
}

// Count counts records
func (g *GormQueryAdapter) Count(count *int64) Query {
	return &GormQueryAdapter{db: g.db.Count(count)}
}

// Distinct adds a DISTINCT clause
func (g *GormQueryAdapter) Distinct(column string) Query {
	return &GormQueryAdapter{db: g.db.Distinct(column)}
}

// Pluck plucks specific columns
func (g *GormQueryAdapter) Pluck(column string, dest interface{}) Query {
	return &GormQueryAdapter{db: g.db.Pluck(column, dest)}
}

// First finds the first record
func (g *GormQueryAdapter) First(dest interface{}, conds ...interface{}) Query {
	return &GormQueryAdapter{db: g.db.First(dest, conds...)}
}
