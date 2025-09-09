package handlers

import (
	"fmt"
	"linn221/Requester/requests"
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"net/http"
	"strconv"
)

// RequestsHandler handles requests related operations
type RequestsHandler struct {
	services *services.ServiceContainer
}

// NewRequestsHandler creates a new RequestsHandler
func NewRequestsHandler(services *services.ServiceContainer) *RequestsHandler {
	return &RequestsHandler{
		services: services,
	}
}

// HandleRequestsList handles GET /requests
func (h *RequestsHandler) HandleRequestsList(w http.ResponseWriter, r *http.Request) error {
	// Parse query parameters
	importJobIDStr := r.URL.Query().Get("import_job_id")
	endpointIDStr := r.URL.Query().Get("endpoint_id")
	orderBy := r.URL.Query().Get("orderBy")
	if orderBy == "" {
		orderBy = r.URL.Query().Get("order_by") // Fallback for old parameter name
	}
	direction := r.URL.Query().Get("direction")
	if direction == "" {
		// Fallback for old asc parameter
		ascStr := r.URL.Query().Get("asc")
		direction = "ASC"
		if ascStr == "false" {
			direction = "DESC"
		}
	}
	search := r.URL.Query().Get("search")

	var requests []requests.MyRequest
	var pageTitle string

	// Check if filtering by endpoint_id
	if endpointIDStr != "" {
		endpointID, err := strconv.ParseUint(endpointIDStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid endpoint ID: %v", err)
		}

		// Fetch requests by endpoint
		requests, err = h.services.RequestService.GetRequestsByEndpointWithOrderAndSearch(r.Context(), uint(endpointID), orderBy, direction, search)
		if err != nil {
			return err
		}
		pageTitle = "Endpoint Requests"
	} else if importJobIDStr != "" {
		importJobID, err := strconv.ParseUint(importJobIDStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid import job ID: %v", err)
		}

		// Fetch requests by import job
		requests, err = h.services.RequestService.GetRequestsByImportJobWithOrderAndSearch(r.Context(), uint(importJobID), orderBy, direction, search)
		if err != nil {
			return err
		}

		// Get import job title
		importJob, err := h.services.ImportJobService.GetImportJobByID(r.Context(), uint(importJobID))
		if err != nil {
			return err
		}
		pageTitle = importJob.Title
	} else {
		return fmt.Errorf("either import_job_id or endpoint_id parameter is required")
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.RequestsList(requests, orderBy, direction, importJobIDStr, endpointIDStr, search, pageTitle).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.RequestsListPage(requests, orderBy, direction, importJobIDStr, endpointIDStr, search, pageTitle).Render(r.Context(), w)
	}
}
