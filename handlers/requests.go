package handlers

import (
	"fmt"
	"linn221/Requester/requests"
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"net/http"
	"strconv"
	"strings"
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
	search := r.URL.Query().Get("search")

	// Parse multi-order parameters
	orders := h.parseMultiOrderParams(r)

	var requests []requests.MyRequest
	var pageTitle string

	// Check if filtering by endpoint_id
	if endpointIDStr != "" && endpointIDStr != "0" {
		endpointID, err := strconv.ParseUint(endpointIDStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid endpoint ID: %v", err)
		}

		// Fetch requests by endpoint with multi-order
		requests, err = h.services.RequestService.GetRequestsByEndpointWithMultiOrderAndSearch(r.Context(), uint(endpointID), orders, search)
		if err != nil {
			return err
		}
		pageTitle = "Endpoint Requests"
	} else if importJobIDStr != "" && importJobIDStr != "0" {
		importJobID, err := strconv.ParseUint(importJobIDStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid import job ID: %v", err)
		}

		// Fetch requests by import job with multi-order
		requests, err = h.services.RequestService.GetRequestsByImportJobWithMultiOrderAndSearch(r.Context(), uint(importJobID), orders, search)
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

	// Create filter state for template
	filterState := h.createFilterState(importJobIDStr, endpointIDStr, search, orders)

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.RequestsList(requests, filterState, pageTitle).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.RequestsListPage(requests, filterState, pageTitle).Render(r.Context(), w)
	}
}

// createFilterState creates the filter state for the template
func (h *RequestsHandler) createFilterState(importJobID, endpointID, search string, orders []services.OrderClause) templates.FilterState {
	// Ensure we have at least 4 order slots, but only show the first one initially
	orderClauses := make([]templates.OrderClause, 4)
	
	// Set the first order (always show this one)
	if len(orders) > 0 {
		orderClauses[0] = templates.OrderClause{
			Column:    orders[0].Column,
			Direction: orders[0].Direction,
		}
	} else {
		orderClauses[0] = templates.OrderClause{
			Column:    "created_at",
			Direction: "desc",
		}
	}
	
	// Set additional orders if they exist
	for i := 1; i < 4; i++ {
		if i < len(orders) {
			orderClauses[i] = templates.OrderClause{
				Column:    orders[i].Column,
				Direction: orders[i].Direction,
			}
		} else {
			orderClauses[i] = templates.OrderClause{
				Column:    "",
				Direction: "desc",
			}
		}
	}

	return templates.FilterState{
		Search:      search,
		ImportJobID: importJobID,
		EndpointID:  endpointID,
		Orders:      orderClauses,
	}
}

// parseMultiOrderParams parses order_0, direction_0, order_1, direction_1, etc. parameters
func (h *RequestsHandler) parseMultiOrderParams(r *http.Request) []services.OrderClause {
	var orders []services.OrderClause

	// Check for multi-order parameters (order_0, direction_0, etc.)
	for i := 0; i < 4; i++ { // Support up to 4 order clauses
		orderParam := r.URL.Query().Get(fmt.Sprintf("order_%d", i))
		directionParam := r.URL.Query().Get(fmt.Sprintf("direction_%d", i))

		if orderParam != "" {
			// Normalize direction
			direction := strings.ToUpper(directionParam)
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}

			orders = append(orders, services.OrderClause{
				Column:    orderParam,
				Direction: direction,
			})
		}
	}

	// If no multi-order parameters found, fall back to single order parameters
	if len(orders) == 0 {
		orderBy := r.URL.Query().Get("orderBy")
		if orderBy == "" {
			orderBy = r.URL.Query().Get("order_by")
		}
		direction := r.URL.Query().Get("direction")
		if direction == "" {
			ascStr := r.URL.Query().Get("asc")
			direction = "ASC"
			if ascStr == "false" {
				direction = "DESC"
			}
		}

		if orderBy != "" {
			orders = append(orders, services.OrderClause{
				Column:    orderBy,
				Direction: direction,
			})
		}
	}

	// Default ordering if no orders specified
	if len(orders) == 0 {
		orders = append(orders, services.OrderClause{
			Column:    "created_at",
			Direction: "DESC",
		})
	}

	return orders
}

// HandleRequestDetail handles GET /requests/{id}
func (h *RequestsHandler) HandleRequestDetail(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid request ID: %v", err)
	}

	// Fetch request
	request, err := h.services.RequestService.GetRequestByID(r.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.RequestDetail(*request).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.RequestDetailPage(*request).Render(r.Context(), w)
	}
}
