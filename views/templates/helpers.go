package templates

import (
	"encoding/json"
	"fmt"
	"net/url"
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

// truncateURL truncates a URL to the specified length and adds ellipsis
func truncateURL(url string, maxLength int) string {
	if len(url) <= maxLength {
		return url
	}
	return url[:maxLength-3] + "..."
}

// ExtractURIWithoutQuery extracts the URI path without query parameters
func ExtractURIWithoutQuery(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil {
		return fullURL
	}
	return u.Path
}

// parseHeaders parses JSON headers into a slice of header maps
func parseHeaders(headersJSON string) []map[string]string {
	if headersJSON == "" {
		return []map[string]string{}
	}

	var headers []map[string]string
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return []map[string]string{}
	}

	return headers
}

// isHTMLContent checks if the content type indicates HTML
func isHTMLContent(headersJSON string) bool {
	contentType := getContentType(headersJSON)
	return strings.Contains(strings.ToLower(contentType), "text/html")
}

// isJSONContent checks if the content type indicates JSON
func isJSONContent(headersJSON string) bool {
	contentType := getContentType(headersJSON)
	return strings.Contains(strings.ToLower(contentType), "application/json")
}

// formatJSON formats JSON string with proper indentation
func formatJSON(jsonStr string) string {
	if jsonStr == "" {
		return ""
	}

	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return jsonStr // Return original if not valid JSON
	}

	formatted, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return jsonStr // Return original if formatting fails
	}

	return string(formatted)
}

// getEndpointIDsJS converts endpoint IDs slice to JavaScript array string
func getEndpointIDsJS(endpointIDs []string) string {
	if len(endpointIDs) == 0 {
		return "[]"
	}

	jsArray := "["
	for i, id := range endpointIDs {
		if i > 0 {
			jsArray += ","
		}
		jsArray += fmt.Sprintf(`"%s"`, id)
	}
	jsArray += "]"
	return jsArray
}

// getMethodsJS converts methods slice to JavaScript array string
func getMethodsJS(methods []string) string {
	if len(methods) == 0 {
		return "[]"
	}

	jsArray := "["
	for i, method := range methods {
		if i > 0 {
			jsArray += ","
		}
		jsArray += fmt.Sprintf(`"%s"`, method)
	}
	jsArray += "]"
	return jsArray
}

// getStatusesJS converts statuses slice to JavaScript array string
func getStatusesJS(statuses []string) string {
	if len(statuses) == 0 {
		return "[]"
	}

	jsArray := "["
	for i, status := range statuses {
		if i > 0 {
			jsArray += ","
		}
		jsArray += fmt.Sprintf(`"%s"`, status)
	}
	jsArray += "]"
	return jsArray
}

// getTypesJS converts types slice to JavaScript array string
func getTypesJS(types []string) string {
	if len(types) == 0 {
		return "[]"
	}

	jsArray := "["
	for i, typeStr := range types {
		if i > 0 {
			jsArray += ","
		}
		jsArray += fmt.Sprintf(`"%s"`, typeStr)
	}
	jsArray += "]"
	return jsArray
}
