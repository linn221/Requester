package main

import (
	"fmt"
	"linn221/Requester/handlers"
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

func MakeAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookiesSecret, err := r.Cookie("secret")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Redirect(w, r, "/secret-required", http.StatusTemporaryRedirect)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if cookiesSecret.Value == secret {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/secret-required", http.StatusTemporaryRedirect)
			}
		})
	}
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic message and stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", rec, debug.Stack())

				// Optional: customize the error response
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func makeDashboardRoutes(app *App, mux *http.ServeMux) {
	// Create handlers
	importJobsHandler := handlers.NewImportJobsHandler(app.services)
	requestsHandler := handlers.NewRequestsHandler(app.services)
	endpointsHandler := handlers.NewEndpointsHandler(app.services)
	programsHandler := handlers.NewProgramsHandler(app.services)

	// Home page - check if it's an HTMX request
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "true" {
			// HTMX request - return just the content
			templates.HomePage().Render(r.Context(), w)
		} else {
			// Direct visit - return full page with layout
			templates.LayoutWithNav("Home - App", templates.HomePage(), "home").Render(r.Context(), w)
		}
	})

	// Import form - check if it's an HTMX request
	mux.HandleFunc("GET /import", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		// Fetch all programs for the dropdown
		programs, err := app.services.ProgramService.GetAllPrograms(r.Context())
		if err != nil {
			return err
		}

		if r.Header.Get("HX-Request") == "true" {
			// HTMX request - return just the form
			return templates.ImportForm(programs).Render(r.Context(), w)
		} else {
			// Direct visit - return full page with layout
			return templates.ImportFormPage(programs).Render(r.Context(), w)
		}
	}))

	mux.HandleFunc("POST /import", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return handleImport(app, w, r)
	}))

	// Programs CRUD routes
	mux.HandleFunc("GET /programs", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramsList(w, r)
	}))
	mux.HandleFunc("GET /programs/create", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramCreate(w, r)
	}))
	mux.HandleFunc("POST /programs", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramStore(w, r)
	}))
	mux.HandleFunc("GET /programs/{id}/edit", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramEdit(w, r)
	}))
	mux.HandleFunc("PUT /programs/{id}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramUpdate(w, r)
	}))
	mux.HandleFunc("DELETE /programs/{id}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return programsHandler.HandleProgramDelete(w, r)
	}))

	// Import jobs list - check if it's an HTMX request
	mux.HandleFunc("GET /import-jobs", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return importJobsHandler.HandleImportJobsList(w, r)
	}))

	// Endpoints list - check if it's an HTMX request
	mux.HandleFunc("GET /endpoints", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return endpointsHandler.HandleEndpointsList(w, r)
	}))

	// Endpoint detail - check if it's an HTMX request
	mux.HandleFunc("GET /endpoints/{id}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return endpointsHandler.HandleEndpointDetail(w, r)
	}))

	// Requests list - check if it's an HTMX request
	mux.HandleFunc("GET /requests", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return requestsHandler.HandleRequestsList(w, r)
	}))

	// Request detail - check if it's an HTMX request
	mux.HandleFunc("GET /requests/detail/{id}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return handleRequestDetail(app, w, r)
	}))
}

func handleImport(app *App, w http.ResponseWriter, r *http.Request) error {
	// Parse form using service
	importReq, err := app.services.FormParser.ParseImportForm(services.NewHTTPRequestAdapter(r))
	if err != nil {
		return err
	}

	// Import HAR using service
	result, err := app.services.ImportService.ImportHAR(r.Context(), *importReq)
	if err != nil {
		return err
	}

	// Create summary (TODO: Get actual stats from service)
	summary := templates.ImportSummary{
		TotalRequests:   result.RequestCount,
		UniqueEndpoints: 0, // TODO: Get from service
		UniqueDomains:   result.UniqueDomains,
		Methods:         make(map[string]int),
		StatusCodes:     make(map[string]int),
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return templates.ImportResult(importReq.Title, result.RequestCount, result.UniqueDomains, summary, result.ImportJobID).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return templates.ImportResultPage(importReq.Title, result.RequestCount, result.UniqueDomains, summary, result.ImportJobID).Render(r.Context(), w)
	}
}

func handleRequestDetail(app *App, w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid request ID: %v", err)
	}

	// Fetch request using service
	request, err := app.services.RequestService.GetRequestByID(r.Context(), uint(id))
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

func (a *App) Serve() {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMiddleware := MakeAuthMiddleware(a.secret)

	// Dashboard routes with authentication
	makeDashboardRoutes(a, authMux)

	// Authentication routes
	mux.HandleFunc("GET /start-session", func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("secret")
		if s != "" {
			// set cookies
			http.SetCookie(w, &http.Cookie{
				Name:     "secret",
				Value:    s,
				MaxAge:   10 * 60 * 60,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
			})

			http.Redirect(w, r, "/dashboard/", http.StatusTemporaryRedirect) // the final trail is important for some reason
			return
		}
	})

	mux.HandleFunc("GET /secret-required", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Secrets required. Open the link from running the app"))
	})

	mux.Handle("/dashboard/", http.StripPrefix("/dashboard", authMiddleware(authMux))) // the final trail is important

	// Print the URL with secret for authentication
	fmt.Printf("http://localhost:%s/start-session?secret=%s\n", a.port, a.secret)

	err := http.ListenAndServe(":"+a.port, RecoveryMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
