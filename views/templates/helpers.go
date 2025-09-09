package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// getMethodBadgeClass returns the appropriate CSS class for HTTP method badges
func getMethodBadgeClass(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "badge bg-success"
	case "POST":
		return "badge bg-primary"
	case "PUT":
		return "badge bg-warning"
	case "PATCH":
		return "badge bg-info"
	case "DELETE":
		return "badge bg-danger"
	case "HEAD":
		return "badge bg-secondary"
	default:
		return "badge bg-dark"
	}
}

// getStatusBadgeClass returns the appropriate CSS class for HTTP status badges
func getStatusBadgeClass(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "badge bg-warning" // Yellow for 2xx
	case status >= 300 && status < 400:
		return "badge bg-success" // Green for 3xx
	case status >= 400 && status < 500:
		return "badge bg-primary" // Blue for 4xx
	case status >= 500:
		return "badge bg-danger" // Red for 5xx
	default:
		return "badge bg-secondary"
	}
}

// formatBytes formats bytes into human readable format
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := int64(bytes) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDate formats a timestamp into a readable date string
func formatDate(timestamp int64) string {
	if timestamp == 0 {
		return "Never"
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// getContentType extracts content-type from response headers JSON
func getContentType(headersJSON string) string {
	if headersJSON == "" {
		return "Unknown"
	}
	
	var headers []map[string]string
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return "Unknown"
	}
	
	for _, header := range headers {
		if strings.ToLower(header["name"]) == "content-type" {
			return header["value"]
		}
	}
	
	return "Unknown"
}
