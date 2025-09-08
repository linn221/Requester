package main

import (
	"fmt"
	"io"
	"linn221/Requester/requests"
	"linn221/Requester/views"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
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
	// Home page - check if it's an HTMX request
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") == "true" {
			// HTMX request - return just the content
			views.HomePage().Render(r.Context(), w)
		} else {
			// Direct visit - return full page with layout
			views.Layout("Home - App", views.HomePage()).Render(r.Context(), w)
		}
	})

	// Import form - check if it's an HTMX request
	mux.HandleFunc("GET /import", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		if r.Header.Get("HX-Request") == "true" {
			// HTMX request - return just the form
			return views.ImportForm().Render(r.Context(), w)
		} else {
			// Direct visit - return full page with layout
			return views.ImportFormPage().Render(r.Context(), w)
		}
	}))

	mux.HandleFunc("POST /import", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return handleImport(app, w, r)
	}))

	// Requests list - check if it's an HTMX request
	mux.HandleFunc("GET /requests/{importJobId}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return handleRequestsList(app, w, r)
	}))

	// Request detail - check if it's an HTMX request
	mux.HandleFunc("GET /requests/detail/{id}", app.HandleMin(func(w http.ResponseWriter, r *http.Request) error {
		return handleRequestDetail(app, w, r)
	}))
}

func handleImport(app *App, w http.ResponseWriter, r *http.Request) error {
	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Get form values
	title := r.FormValue("title")
	if title == "" {
		return fmt.Errorf("title is required")
	}

	ignoredHeadersText := r.FormValue("ignoredHeaders")
	ignoredHeaders := strings.Fields(strings.ReplaceAll(ignoredHeadersText, "\n", " "))

	// Get uploaded file
	file, header, err := r.FormFile("harfile")
	if err != nil {
		return fmt.Errorf("failed to get uploaded file: %v", err)
	}
	defer file.Close()

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".har") {
		return fmt.Errorf("file must be a .har file")
	}

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Create transaction with context
	tx := app.db.WithContext(r.Context()).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create ImportJob record
	importJob := requests.ImportJob{
		Title:          title,
		IgnoredHeaders: strings.Join(ignoredHeaders, ","),
	}

	// Save ImportJob to database
	if err := tx.Create(&importJob).Error; err != nil {
		return fmt.Errorf("failed to create import job: %v", err)
	}

	// Create resHashFunc that uses ignored headers
	resHashFunc := func(my *requests.TempMyRequest) (string, string) {
		// Request text with filtered headers
		reqText := my.URL + " " + my.Method + " " + my.ReqBody + " " + my.ReqHeaders.EchoFilter(ignoredHeaders...)

		// Response text with filtered headers
		respText := fmt.Sprintf("%d %d %s %s",
			my.ResStatus, my.RespSize, my.ResBody, my.ResHeaders.EchoFilter(ignoredHeaders...),
		)
		return reqText, respText
	}

	// Parse HAR file
	tempResults, err := requests.ParseHAR(fileContent, resHashFunc)
	if err != nil {
		return fmt.Errorf("failed to parse HAR file: %v", err)
	}

	// Convert TempMyRequest to MyRequest and save to database
	var dbResults []requests.MyRequest
	for _, tempReq := range tempResults {
		dbReq, err := tempReq.ToMyRequest(importJob.ID)
		if err != nil {
			return fmt.Errorf("failed to convert request to database format: %v", err)
		}
		dbResults = append(dbResults, *dbReq)
	}

	// Save all requests to database in batch
	if err := tx.CreateInBatches(dbResults, 100).Error; err != nil {
		return fmt.Errorf("failed to save requests to database: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate summary
	summary := generateImportSummary(tempResults, title)

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return views.ImportResult(title, len(tempResults), countUniqueDomains(tempResults), summary, importJob.ID).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return views.ImportResultPage(title, len(tempResults), countUniqueDomains(tempResults), summary, importJob.ID).Render(r.Context(), w)
	}
}

func generateImportSummary(results []requests.TempMyRequest, title string) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("HAR Import Summary for: %s\n", title))
	summary.WriteString(fmt.Sprintf("Total Requests: %d\n", len(results)))
	summary.WriteString(fmt.Sprintf("Unique Domains: %d\n", countUniqueDomains(results)))
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

func countUniqueDomains(results []requests.TempMyRequest) int {
	domains := make(map[string]bool)
	for _, req := range results {
		domains[req.Domain] = true
	}
	return len(domains)
}

func handleRequestsList(app *App, w http.ResponseWriter, r *http.Request) error {
	// Extract importJobId from URL
	importJobIdStr := r.PathValue("importJobId")
	importJobId, err := strconv.ParseUint(importJobIdStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid import job ID: %v", err)
	}

	// Create transaction with context
	tx := app.db.WithContext(r.Context())

	// Fetch requests for the import job
	var requests []requests.MyRequest
	if err := tx.Where("import_job_id = ?", importJobId).Order("sequence ASC").Find(&requests).Error; err != nil {
		return fmt.Errorf("failed to fetch requests: %v", err)
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return views.RequestsList(requests).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return views.RequestsListPage(requests).Render(r.Context(), w)
	}
}

func handleRequestDetail(app *App, w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid request ID: %v", err)
	}

	// Create transaction with context
	tx := app.db.WithContext(r.Context())

	// Fetch the specific request
	var request requests.MyRequest
	if err := tx.First(&request, id).Error; err != nil {
		return fmt.Errorf("failed to fetch request: %v", err)
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// HTMX request - return just the content
		return views.RequestDetail(request).Render(r.Context(), w)
	} else {
		// Direct visit - return full page with layout
		return views.RequestDetailPage(request).Render(r.Context(), w)
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
