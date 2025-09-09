package handlers

import (
	"fmt"
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
	orderBy := r.URL.Query().Get("order_by")
	
	if importJobIDStr == "" {
		return fmt.Errorf("import_job_id parameter is required")
	}
	
	importJobID, err := strconv.ParseUint(importJobIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid import job ID: %v", err)
	}

	// Fetch requests using service with ordering
	requests, err := h.services.RequestService.GetRequestsByImportJobWithOrder(r.Context(), uint(importJobID), orderBy)
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.RequestsList(requests, orderBy).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.RequestsListPage(requests, orderBy).Render(r.Context(), w)
	}
}
