package services

import (
	"context"
	"fmt"
	"linn221/Requester/requests"
)

// ImportJobService handles import job-related operations
type ImportJobService struct {
	db Database
}

// NewImportJobService creates a new ImportJobService
func NewImportJobService(db Database) *ImportJobService {
	return &ImportJobService{db: db}
}

// GetAllImportJobs fetches all import jobs
func (s *ImportJobService) GetAllImportJobs(ctx context.Context) ([]requests.ImportJob, error) {
	// Create transaction with context
	tx := s.db.WithContext(ctx)

	// Fetch all import jobs
	var importJobs []requests.ImportJob
	if err := tx.Where("1 = 1").Order("created_at DESC").Find(&importJobs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch import jobs: %v", err)
	}

	return importJobs, nil
}
