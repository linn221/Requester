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
func (s *RequestService) GetRequestsByImportJobWithOrder(ctx context.Context, importJobID uint, orderBy, direction string) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest

	// Build the order clause
	orderClause := s.buildOrderClause(orderBy, direction)

	if err := s.db.WithContext(ctx).Where("import_job_id = ?", importJobID).Order(orderClause).Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for job %d: %v", importJobID, err)
	}
	return reqs, nil
}

// GetRequestsByImportJobWithOrderAndSearch fetches requests for a specific import job with custom ordering and search
func (s *RequestService) GetRequestsByImportJobWithOrderAndSearch(ctx context.Context, importJobID uint, orderBy, direction, search string) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest

	// Build the order clause
	orderClause := s.buildOrderClause(orderBy, direction)

	query := s.db.WithContext(ctx).Where("import_job_id = ?", importJobID)

	// Add search condition if search term is provided
	if search != "" {
		query = query.Where("raw LIKE ? OR raw1 LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order(orderClause).Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for job %d: %v", importJobID, err)
	}
	return reqs, nil
}

// GetRequestsByEndpointWithOrder fetches requests for a specific endpoint with custom ordering
func (s *RequestService) GetRequestsByEndpointWithOrder(ctx context.Context, endpointID uint, orderBy, direction string) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest

	// Build the order clause
	orderClause := s.buildOrderClause(orderBy, direction)

	if err := s.db.WithContext(ctx).Where("endpoint_id = ?", endpointID).Order(orderClause).Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for endpoint %d: %v", endpointID, err)
	}
	return reqs, nil
}

// GetRequestsByEndpointWithOrderAndSearch fetches requests for a specific endpoint with custom ordering and search
func (s *RequestService) GetRequestsByEndpointWithOrderAndSearch(ctx context.Context, endpointID uint, orderBy, direction, search string) ([]requests.MyRequest, error) {
	var reqs []requests.MyRequest

	// Build the order clause
	orderClause := s.buildOrderClause(orderBy, direction)

	query := s.db.WithContext(ctx).Where("endpoint_id = ?", endpointID)

	// Add search condition if search term is provided
	if search != "" {
		query = query.Where("raw LIKE ? OR raw1 LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order(orderClause).Find(&reqs).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch requests for endpoint %d: %v", endpointID, err)
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

// buildOrderClause builds the ORDER BY clause based on the orderBy and direction parameters
func (s *RequestService) buildOrderClause(orderBy, direction string) string {
	if orderBy == "" {
		return "sequence ASC"
	}

	// Validate and normalize direction
	if direction == "" {
		direction = "ASC"
	}
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}

	// Map column names to database fields
	switch orderBy {
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
	case "type":
		// For content type, we'll sort by response headers JSON
		// This is a simplified approach - in practice, you might want to extract content-type to a separate field
		return fmt.Sprintf("res_headers %s", direction)
	case "created":
		return fmt.Sprintf("created_at %s", direction)
	default:
		return "sequence ASC"
	}
}
