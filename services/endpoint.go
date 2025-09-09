package services

import (
	"context"
	"fmt"
	"linn221/Requester/requests"
)

// EndpointService handles endpoint-related operations
type EndpointService struct {
	db Database
}

// NewEndpointService creates a new EndpointService
func NewEndpointService(db Database) *EndpointService {
	return &EndpointService{db: db}
}

// FindOrCreateEndpoint finds an existing endpoint or creates a new one
func (s *EndpointService) FindOrCreateEndpoint(ctx context.Context, method, domain, uri string) (*requests.Endpoint, error) {
	// First, try to find existing endpoint
	var endpoint requests.Endpoint
	err := s.db.WithContext(ctx).Where("method = ? AND domain = ? AND uri = ?", method, domain, uri).First(&endpoint).Error()

	if err == nil {
		// Endpoint found, return it
		return &endpoint, nil
	}

	// Endpoint not found, create new one
	endpointType := requests.DetermineEndpointType(uri, method)
	newEndpoint := &requests.Endpoint{
		Method:       method,
		Domain:       domain,
		URI:          uri,
		EndpointType: endpointType,
		Notes:        "",
	}

	if err := s.db.WithContext(ctx).Create(newEndpoint).Error(); err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %v", err)
	}

	return newEndpoint, nil
}

// GetAllEndpoints fetches all endpoints
func (s *EndpointService) GetAllEndpoints(ctx context.Context) ([]requests.Endpoint, error) {
	var endpoints []requests.Endpoint
	if err := s.db.WithContext(ctx).Order("domain ASC, method ASC, uri ASC").Find(&endpoints).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch endpoints: %v", err)
	}
	return endpoints, nil
}

// GetEndpointByID fetches a single endpoint by ID
func (s *EndpointService) GetEndpointByID(ctx context.Context, id uint) (*requests.Endpoint, error) {
	var endpoint requests.Endpoint
	if err := s.db.WithContext(ctx).First(&endpoint, id).Error(); err != nil {
		return nil, fmt.Errorf("failed to fetch endpoint with ID %d: %v", id, err)
	}
	return &endpoint, nil
}

// GetEndpointStats fetches statistics for an endpoint
func (s *EndpointService) GetEndpointStats(ctx context.Context, endpointID uint) (map[string]interface{}, error) {
	var count int64
	var avgLatency float64
	var statusCounts []struct {
		Status int
		Count  int64
	}

	// Get total request count
	if err := s.db.WithContext(ctx).Model(&requests.MyRequest{}).Where("endpoint_id = ?", endpointID).Count(&count).Error(); err != nil {
		return nil, fmt.Errorf("failed to count requests: %v", err)
	}

	// Get average latency
	if err := s.db.WithContext(ctx).Model(&requests.MyRequest{}).Where("endpoint_id = ?", endpointID).Select("AVG(latency_ms)").Scan(&avgLatency).Error(); err != nil {
		return nil, fmt.Errorf("failed to calculate average latency: %v", err)
	}

	// Get status code distribution
	if err := s.db.WithContext(ctx).Model(&requests.MyRequest{}).
		Select("res_status as status, COUNT(*) as count").
		Where("endpoint_id = ?", endpointID).
		Group("res_status").
		Scan(&statusCounts).Error(); err != nil {
		return nil, fmt.Errorf("failed to get status distribution: %v", err)
	}

	stats := map[string]interface{}{
		"total_requests": count,
		"avg_latency":    avgLatency,
		"status_counts":  statusCounts,
	}

	return stats, nil
}
