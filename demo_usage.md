# Request UI Templates Usage

## Overview
I've created two new UI templates for listing and viewing requests:

1. **RequestsList** - A searchable list of requests with cards
2. **RequestDetail** - Detailed view of a single request

## Features

### RequestsList Template
- **Search functionality**: Real-time search across URL, method, status, domain
- **Request cards**: Display URL, method, status, size, latency
- **Click to view**: Click any card to navigate to detail page
- **Responsive design**: Works on mobile and desktop
- **Alpine.js powered**: Reactive search and filtering

### RequestDetail Template
- **Complete request info**: All MyRequest fields displayed
- **Action buttons**: Copy as cURL, Open in VS Code, Add Notes
- **Copy functionality**: Copy headers, bodies, and other data
- **Notes form**: Collapsible form for adding notes
- **Organized layout**: Information grouped logically

## Usage in Routes

```go
// In your routes.go file
func handleRequestsList(app *App, w http.ResponseWriter, r *http.Request) error {
    // Fetch requests from database
    var requests []requests.MyRequest
    if err := app.db.WithContext(r.Context()).Find(&requests).Error; err != nil {
        return err
    }
    
    // Convert to template format
    var listItems []views.RequestListItem
    for _, req := range requests {
        listItems = append(listItems, views.RequestListItem{
            ID:      req.ID,
            URL:     req.URL,
            Method:  req.Method,
            Domain:  req.Domain,
            Status:  req.ResStatus,
            Size:    req.RespSize,
            Latency: req.LatencyMs,
        })
    }
    
    return views.RequestsList(listItems).Render(r.Context(), w)
}

func handleRequestDetail(app *App, w http.ResponseWriter, r *http.Request) error {
    // Extract ID from URL
    id := chi.URLParam(r, "id")
    
    // Fetch request from database
    var req requests.MyRequest
    if err := app.db.WithContext(r.Context()).First(&req, id).Error; err != nil {
        return err
    }
    
    // Convert to template format
    detailData := views.RequestDetailData{
        ID:            req.ID,
        ImportJobID:   req.ImportJobID,
        Sequence:      req.Sequence,
        URL:           req.URL,
        Method:        req.Method,
        Domain:        req.Domain,
        ReqHeaders:    req.ReqHeaders,
        ReqBody:       req.ReqBody,
        Status:        req.ResStatus,
        ResHeaders:    req.ResHeaders,
        ResBody:       req.ResBody,
        Size:          req.RespSize,
        Latency:       req.LatencyMs,
        RequestTime:   req.RequestTime,
        ReqHash1:      req.ReqHash1,
        ReqHash:       req.ReqHash,
        ResHash:       req.ResHash,
        ResBodyHash:   req.ResBodyHash,
        CreatedAt:     req.CreatedAt,
        UpdatedAt:     req.UpdatedAt,
    }
    
    return views.RequestDetail(detailData).Render(r.Context(), w)
}
```

## Styling
- Uses Bootstrap 5 for responsive design
- Custom CSS for hover effects and transitions
- Bootstrap Icons for consistent iconography
- Alpine.js for reactivity without complexity

## Dependencies Added
- Alpine.js 3.x for reactivity
- Bootstrap Icons for icons
- Custom CSS for enhanced UX

The templates are ready to use and will provide a clean, user-friendly interface for managing requests.
