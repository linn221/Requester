package services

import (
	"context"
	"fmt"
	"linn221/Requester/requests"
	"strings"
)

// ImportService handles HAR file import operations
type ImportService struct {
	db              Database
	endpointService *EndpointService
}

// NewImportService creates a new ImportService
func NewImportService(db Database, endpointService *EndpointService) *ImportService {
	return &ImportService{
		db:              db,
		endpointService: endpointService,
	}
}

// ImportHAR processes a HAR file import
func (s *ImportService) ImportHAR(ctx context.Context, req ImportRequest) (*ImportResult, error) {
	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(req.Filename), ".har") {
		return nil, fmt.Errorf("file must be a .har file")
	}

	// Create transaction with context
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error() != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create ImportJob record
	importJob := requests.ImportJob{
		Title:          req.Title,
		IgnoredHeaders: strings.Join(req.IgnoredHeaders, ","),
	}

	// Save ImportJob to database
	if err := tx.Create(&importJob).Error(); err != nil {
		return nil, fmt.Errorf("failed to create import job: %v", err)
	}

	// Create resHashFunc that uses ignored headers
	resHashFunc := func(my *requests.TempMyRequest) (string, string) {
		// Request text with filtered headers
		reqText := my.URL + " " + my.Method + " " + my.ReqBody + " " + my.ReqHeaders.EchoFilter(req.IgnoredHeaders...)

		// Response text with filtered headers
		respText := fmt.Sprintf("%d %d %s %s",
			my.ResStatus, my.RespSize, my.ResBody, my.ResHeaders.EchoFilter(req.IgnoredHeaders...),
		)
		return reqText, respText
	}

	// Parse HAR file
	tempResults, err := requests.ParseHAR(req.FileContent, resHashFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HAR file: %v", err)
	}

	// Convert TempMyRequest to MyRequest and save to database
	var dbResults []requests.MyRequest
	for _, tempReq := range tempResults {
		// Extract URI without query parameters
		uri := requests.ExtractURIWithoutQuery(tempReq.URL)

		// Find or create endpoint
		endpoint, err := s.endpointService.FindOrCreateEndpoint(ctx, tempReq.Method, tempReq.Domain, uri)
		if err != nil {
			return nil, fmt.Errorf("failed to find or create endpoint: %v", err)
		}

		dbReq, err := tempReq.ToMyRequest(importJob.ID, endpoint.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert request to database format: %v", err)
		}
		dbResults = append(dbResults, *dbReq)
	}

	// Save all requests to database in batch
	if err := tx.CreateInBatches(dbResults, 100).Error(); err != nil {
		return nil, fmt.Errorf("failed to save requests to database: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate summary
	summary := GenerateImportSummary(tempResults, req.Title)

	return &ImportResult{
		ImportJobID:   importJob.ID,
		RequestCount:  len(tempResults),
		UniqueDomains: CountUniqueDomains(tempResults),
		Summary:       summary,
	}, nil
}

// GenerateImportSummary creates a summary of the import results
func GenerateImportSummary(results []requests.TempMyRequest, title string) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("HAR Import Summary for: %s\n", title))
	summary.WriteString(fmt.Sprintf("Total Requests: %d\n", len(results)))
	summary.WriteString(fmt.Sprintf("Unique Domains: %d\n", CountUniqueDomains(results)))
	summary.WriteString("\nDomain Breakdown:\n")

	domainCounts := make(map[string]int)
	for _, req := range results {
		domainCounts[req.Domain]++
	}

	for domain, count := range domainCounts {
		summary.WriteString(fmt.Sprintf("  %s: %d requests\n", domain, count))
	}

	summary.WriteString("\nMethod Breakdown:\n")
	methodCounts := make(map[string]int)
	for _, req := range results {
		methodCounts[req.Method]++
	}

	for method, count := range methodCounts {
		summary.WriteString(fmt.Sprintf("  %s: %d requests\n", method, count))
	}

	summary.WriteString("\nStatus Code Breakdown:\n")
	statusCounts := make(map[int]int)
	for _, req := range results {
		statusCounts[req.ResStatus]++
	}

	for status, count := range statusCounts {
		summary.WriteString(fmt.Sprintf("  %d: %d responses\n", status, count))
	}

	return summary.String()
}

// CountUniqueDomains counts unique domains in the results
func CountUniqueDomains(results []requests.TempMyRequest) int {
	domains := make(map[string]bool)
	for _, req := range results {
		domains[req.Domain] = true
	}
	return len(domains)
}
