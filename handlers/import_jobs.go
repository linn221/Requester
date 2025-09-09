package handlers

import (
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"net/http"
)

// ImportJobsHandler handles import jobs related requests
type ImportJobsHandler struct {
	services *services.ServiceContainer
}

// NewImportJobsHandler creates a new ImportJobsHandler
func NewImportJobsHandler(services *services.ServiceContainer) *ImportJobsHandler {
	return &ImportJobsHandler{
		services: services,
	}
}

// HandleImportJobsList handles GET /import-jobs
func (h *ImportJobsHandler) HandleImportJobsList(w http.ResponseWriter, r *http.Request) error {
	// Fetch all import jobs
	importJobs, err := h.services.ImportJobService.GetAllImportJobs(r.Context())
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.ImportJobsList(importJobs).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.ImportJobsListPage(importJobs).Render(r.Context(), w)
	}
}
