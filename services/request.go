package services

import (
	"context"
	"fmt"
	"linn221/Requester/requests"
	"strings"
)

// RequestService handles request-related operations
type RequestService struct {
	db Database
}

// NewRequestService creates a new RequestService
func NewRequestService(db Database) *RequestService {
	return &RequestService{db: db}
}

// GetRequestsByImportJob fetches requests for a specific import job
func (s *RequestService) GetRequestsByImportJob(ctx context.Context, importJobID uint) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest
	if err := s.db.WithContext(ctx).Where("import_job_id = ?", importJobID).Order("sequence ASC").Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for job %d: %v", importJobID, err)
	}
	return reqs, nil
}

// GetRequestsByImportJobWithOrder fetches requests for a specific import job with custom ordering
func (s *RequestService) GetRequestsByImportJobWithOrder(ctx context.Context, importJobID uint, orderBy string) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest
	
	// Build the order clause
	orderClause := s.buildOrderClause(orderBy)
	
	if err := s.db.WithContext(ctx).Where("import_job_id = ?", importJobID).Order(orderClause).Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for job %d: %v", importJobID, err)
	}
	return reqs, nil
}

// GetRequestByID fetches a single request by ID
func (s *RequestService) GetRequestByID(ctx context.Context, id uint) (*requests.MyRequest, error) {
	var request requests.MyRequest
	if err := s.db.WithContext(ctx).First(&request, id).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch request with ID %d: %v", id, err)
	}
	return &request, nil
}

// buildOrderClause builds the ORDER BY clause based on the orderBy parameter
func (s *RequestService) buildOrderClause(orderBy string) string {
	if orderBy == "" {
		return "sequence ASC"
	}
	
	// Parse orderBy parameter (format: "column:direction" or just "column")
	parts := strings.Split(orderBy, ":")
	column := parts[0]
	direction := "ASC"
	
	if len(parts) > 1 {
		direction = strings.ToUpper(parts[1])
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}
	}
	
	// Map column names to database fields
	switch column {
	case "sequence_number":
		return fmt.Sprintf("sequence %s", direction)
	case "method":
		return fmt.Sprintf("method %s", direction)
	case "status":
		return fmt.Sprintf("res_status %s", direction)
	case "url":
		return fmt.Sprintf("url %s", direction)
	case "domain":
		return fmt.Sprintf("domain %s", direction)
	case "size":
		return fmt.Sprintf("resp_size %s", direction)
	case "latency":
		return fmt.Sprintf("latency_ms %s", direction)
	case "created":
		return fmt.Sprintf("created_at %s", direction)
	default:
		return "sequence ASC"
	}
}
