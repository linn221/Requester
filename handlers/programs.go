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

// ProgramsHandler handles program related requests
type ProgramsHandler struct {
	services *services.ServiceContainer
}

// NewProgramsHandler creates a new ProgramsHandler
func NewProgramsHandler(services *services.ServiceContainer) *ProgramsHandler {
	return &ProgramsHandler{
		services: services,
	}
}

// HandleProgramsList handles GET /programs
func (h *ProgramsHandler) HandleProgramsList(w http.ResponseWriter, r *http.Request) error {
	// Fetch all programs
	programs, err := h.services.ProgramService.GetAllPrograms(r.Context())
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.ProgramsList(programs).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.ProgramsListPage(programs).Render(r.Context(), w)
	}
}

// HandleProgramCreate handles GET /programs/create
func (h *ProgramsHandler) HandleProgramCreate(w http.ResponseWriter, r *http.Request) error {
	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the form
		return templates.ProgramForm(nil).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.ProgramFormPage(nil).Render(r.Context(), w)
	}
}

// HandleProgramStore handles POST /programs
func (h *ProgramsHandler) HandleProgramStore(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Create program from form data
	program := &requests.Program{
		Name:    strings.TrimSpace(r.FormValue("name")),
		URL:     strings.TrimSpace(r.FormValue("url")),
		Notes:   strings.TrimSpace(r.FormValue("notes")),
		Scope:   strings.TrimSpace(r.FormValue("scope")),
		Domains: strings.TrimSpace(r.FormValue("domains")),
	}

	// Validate required fields
	if program.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Create the program
	if err := h.services.ProgramService.CreateProgram(r.Context(), program); err != nil {
		return err
	}

	// Redirect to programs list
	http.Redirect(w, r, "/dashboard/programs", http.StatusSeeOther)
	return nil
}

// HandleProgramEdit handles GET /programs/{id}/edit
func (h *ProgramsHandler) HandleProgramEdit(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid program ID: %v", err)
	}

	// Fetch program
	program, err := h.services.ProgramService.GetProgramByID(r.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the form
		return templates.ProgramForm(program).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.ProgramFormPage(program).Render(r.Context(), w)
	}
}

// HandleProgramUpdate handles PUT /programs/{id}
func (h *ProgramsHandler) HandleProgramUpdate(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid program ID: %v", err)
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Create program from form data
	program := &requests.Program{
		ID:      uint(id),
		Name:    strings.TrimSpace(r.FormValue("name")),
		URL:     strings.TrimSpace(r.FormValue("url")),
		Notes:   strings.TrimSpace(r.FormValue("notes")),
		Scope:   strings.TrimSpace(r.FormValue("scope")),
		Domains: strings.TrimSpace(r.FormValue("domains")),
	}

	// Validate required fields
	if program.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Update the program
	if err := h.services.ProgramService.UpdateProgram(r.Context(), program); err != nil {
		return err
	}

	// Redirect to programs list
	http.Redirect(w, r, "/dashboard/programs", http.StatusSeeOther)
	return nil
}

// HandleProgramDelete handles DELETE /programs/{id}
func (h *ProgramsHandler) HandleProgramDelete(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid program ID: %v", err)
	}

	// Delete the program
	if err := h.services.ProgramService.DeleteProgram(r.Context(), uint(id)); err != nil {
		return err
	}

	// Redirect to programs list
	http.Redirect(w, r, "/dashboard/programs", http.StatusSeeOther)
	return nil
}
