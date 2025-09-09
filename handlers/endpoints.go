package handlers

import (
	"fmt"
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"net/http"
	"strconv"
)

// EndpointsHandler handles endpoint related operations
type EndpointsHandler struct {
	services *services.ServiceContainer
}

// NewEndpointsHandler creates a new EndpointsHandler
func NewEndpointsHandler(services *services.ServiceContainer) *EndpointsHandler {
	return &EndpointsHandler{
		services: services,
	}
}

// HandleEndpointsList handles GET /endpoints
func (h *EndpointsHandler) HandleEndpointsList(w http.ResponseWriter, r *http.Request) error {
	// Fetch all endpoints
	endpoints, err := h.services.EndpointService.GetAllEndpoints(r.Context())
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.EndpointsList(endpoints).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.EndpointsListPage(endpoints).Render(r.Context(), w)
	}
}

// HandleEndpointDetail handles GET /endpoints/{id}
func (h *EndpointsHandler) HandleEndpointDetail(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid endpoint ID: %v", err)
	}

	// Fetch endpoint
	endpoint, err := h.services.EndpointService.GetEndpointByID(r.Context(), uint(id))
	if err != nil {
		return err
	}

	// Fetch endpoint statistics
	stats, err := h.services.EndpointService.GetEndpointStats(r.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.EndpointDetail(*endpoint, stats).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.EndpointDetailPage(*endpoint, stats).Render(r.Context(), w)
	}
}
