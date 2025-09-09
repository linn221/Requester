# Services Package

This package contains the business logic for the Requester application, organized into separate services for better maintainability and reusability.

## Structure

### Core Services
- **`import.go`** - Handles HAR file import operations
- **`request.go`** - Manages request-related database operations
- **`parser.go`** - Parses HTTP form data

### Infrastructure
- **`interfaces.go`** - Defines interfaces for dependency injection
- **`database.go`** - GORM database adapter
- **`http_adapter.go`** - HTTP request adapter
- **`container.go`** - Service container for dependency management

## Key Benefits

1. **Separation of Concerns**: Business logic is separated from HTTP handling
2. **Reusability**: Services can be used from different packages (web, CLI, etc.)
3. **Testability**: Each service can be unit tested independently
4. **Maintainability**: Clear organization makes code easier to understand and modify
5. **Dependency Injection**: Interfaces allow for easy mocking and testing

## Usage

### In Web Application
```go
// Services are injected into the App struct
app := NewApp(db, "8080", secret)

// Use services in handlers
importReq, err := app.services.FormParser.ParseImportForm(services.NewHTTPRequestAdapter(r))
result, err := app.services.ImportService.ImportHAR(r.Context(), *importReq)
```

### In CLI Application (Future)
```go
// Services can be used directly
importService := services.NewImportService(database)
requestService := services.NewRequestService(database)

// Use services for CLI operations
result, err := importService.ImportHAR(ctx, importRequest)
requests, err := requestService.GetRequestsByImportJob(ctx, importJobID)
```

## Service Responsibilities

### ImportService
- Validates HAR file format
- Creates database transactions
- Parses HAR files using the requests package
- Converts temporary request objects to database models
- Generates import summaries and statistics

### RequestService
- Fetches requests by import job ID
- Retrieves individual requests by ID
- Handles database queries with proper context

### FormParser
- Parses multipart form data
- Validates required fields
- Reads uploaded file content
- Converts form data to service models

## Database Adapter

The `GormDatabaseAdapter` provides a clean interface between GORM and our services, allowing for:
- Easy testing with mock databases
- Consistent error handling
- Context propagation
- Transaction management

## HTTP Adapter

The `HTTPRequestAdapter` adapts Go's standard `*http.Request` to our service interfaces, enabling:
- Clean separation between HTTP and business logic
- Easy testing of form parsing
- Consistent interface across different HTTP frameworks
