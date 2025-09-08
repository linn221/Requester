package requests

import (
	"encoding/json"
	"fmt"
	"linn221/Requester/utils"
	"net/url"
	"strings"
)

type ImportJob struct {
	ID              uint   `gorm:"primaryKey"`
	Title           string `gorm:"not null"`
	IgnoredHeaders  string `gorm:"type:text"` // Store as JSON string
	CreatedAt       int64  `gorm:"autoCreateTime"`
	UpdatedAt       int64  `gorm:"autoUpdateTime"`
	
	// One-to-many relationship
	Requests []MyRequest `gorm:"foreignKey:ImportJobID"`
}

type MyRequest struct {
	ID            uint   `gorm:"primaryKey"`
	ImportJobID   uint   `gorm:"not null;index"` // Foreign key to ImportJob
	Sequence      int    `gorm:"not null"`
	URL           string `gorm:"type:text;not null"`
	Method        string `gorm:"size:10;not null"`
	Domain        string `gorm:"size:255;not null"`

	ReqHeaders    string `gorm:"type:text"` // Store as JSON string
	ReqBody       string `gorm:"type:longtext"`

	ResStatus     int    `gorm:"not null"`
	ResHeaders    string `gorm:"type:text"` // Store as JSON string
	ResBody       string `gorm:"type:longtext"`
	RespSize      int    `gorm:"not null"`
	LatencyMs     int64  `gorm:"not null"`

	RequestTime   string `gorm:"size:50"`
	// hashes
	ReqHash1      string `gorm:"size:64;index"` // hash raw request
	ReqHash       string `gorm:"size:64;index"`
	ResHash       string `gorm:"size:64;index"`
	ResBodyHash   string `gorm:"size:64;index"`
	
	CreatedAt     int64  `gorm:"autoCreateTime"`
	UpdatedAt     int64  `gorm:"autoUpdateTime"`
}

// Temporary struct for parsing HAR files (with HeaderSlice fields)
type TempMyRequest struct {
	Sequence    int
	URL         string
	Method      string
	Domain      string
	ReqHeaders  HeaderSlice
	ReqBody     string
	ResStatus   int
	ResHeaders  HeaderSlice
	ResBody     string
	RespSize    int
	LatencyMs   int64
	RequestTime string
	// hashes
	ReqHash1    string
	ReqHash     string
	ResHash     string
	ResBodyHash string
}

type HeaderSlice []Header

func (hs HeaderSlice) EchoAll() string {
	var result string
	for _, h := range hs {
		result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
	}
	return result
}

func (hs HeaderSlice) EchoMatcher(headerNames ...string) string {
	matchMap := make(map[string]struct{}, len(headerNames))
	for _, hname := range headerNames {
		matchMap[strings.ToLower(hname)] = struct{}{}
	}

	var result string
	for _, h := range hs {
		if _, match := matchMap[strings.ToLower(h.Name)]; match {
			result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
		}
	}
	return result
}

func (hs HeaderSlice) EchoFilter(headerNames ...string) string {
	filterMap := make(map[string]struct{}, len(headerNames))
	for _, hname := range headerNames {
		filterMap[strings.ToLower(hname)] = struct{}{}
	}

	var result string
	for _, h := range hs {
		if _, filter := filterMap[strings.ToLower(h.Name)]; !filter {
			result += fmt.Sprintf("%s: %s\n", h.Name, h.Value)
		}
	}
	return result
}

// Convert HeaderSlice to JSON string for database storage
func (hs HeaderSlice) ToJSON() (string, error) {
	data, err := json.Marshal(hs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Convert JSON string to HeaderSlice from database
func HeaderSliceFromJSON(jsonStr string) (HeaderSlice, error) {
	var hs HeaderSlice
	err := json.Unmarshal([]byte(jsonStr), &hs)
	return hs, err
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s\n", h.Name, h.Value)
}

func (r TempMyRequest) requestText() string {
	raw := r.Method + " " + r.URL + " " + r.ReqBody + " " + r.ReqHeaders.EchoAll()
	return raw
}

// Convert TempMyRequest to MyRequest for database storage
func (temp *TempMyRequest) ToMyRequest(importJobID uint) (*MyRequest, error) {
	// Convert HeaderSlice to JSON strings
	reqHeadersJSON, err := temp.ReqHeaders.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert request headers to JSON: %v", err)
	}
	
	resHeadersJSON, err := temp.ResHeaders.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert response headers to JSON: %v", err)
	}
	
	return &MyRequest{
		ImportJobID: importJobID,
		Sequence:    temp.Sequence,
		URL:         temp.URL,
		Method:      temp.Method,
		Domain:      temp.Domain,
		ReqHeaders:  reqHeadersJSON,
		ReqBody:     temp.ReqBody,
		ResStatus:   temp.ResStatus,
		ResHeaders:  resHeadersJSON,
		ResBody:     temp.ResBody,
		RespSize:    temp.RespSize,
		LatencyMs:   temp.LatencyMs,
		RequestTime: temp.RequestTime,
		ReqHash1:    temp.ReqHash1,
		ReqHash:     temp.ReqHash,
		ResHash:     temp.ResHash,
		ResBodyHash: temp.ResBodyHash,
	}, nil
}

func ParseHAR(bs []byte, resHashFunc func(*TempMyRequest) (string, string)) ([]TempMyRequest, error) {
	var har HAR
	if err := json.Unmarshal(bs, &har); err != nil {
		return nil, err
	}

	var result []TempMyRequest

	for i, entry := range har.Log.Entries {
		reqHeaders := make([]Header, 0, len(entry.Request.Headers))
		for _, h := range entry.Request.Headers {
			reqHeaders = append(reqHeaders, Header{Name: h.Name, Value: h.Value})
		}

		resHeaders := make([]Header, 0, len(entry.Response.Headers))
		for _, h := range entry.Response.Headers {
			resHeaders = append(resHeaders, Header{Name: h.Name, Value: h.Value})
		}

		u, err := url.Parse(entry.Request.URL)
		domain := ""
		if err == nil {
			domain = u.Hostname()
		}

		resBody := entry.Response.Content.Text
		// Decode base64 if needed
		// if strings.ToLower(entry.Response.Content.Encoding) == "base64" {
		// 	decoded, err := decodeBase64(resBody)
		// 	if err == nil {
		// 		resBody = decoded
		// 	}
		// }

		my := TempMyRequest{
			Sequence:    i + 1,
			URL:         entry.Request.URL,
			Domain:      domain,
			ReqHeaders:  reqHeaders,
			ReqBody:     entry.Request.PostData.Text,
			ResHeaders:  resHeaders,
			ResStatus:   entry.Response.Status,
			ResBody:     resBody,
			RespSize:    len(resBody),
			LatencyMs:   int64(entry.Time),
			RequestTime: entry.StartedDateTime,
			Method:      entry.Request.Method,
		}

		requestText, responseText := resHashFunc(&my)

		my.ReqHash = utils.HashString(requestText)
		my.ResHash = utils.HashString(responseText)
		my.ResBodyHash = utils.HashString(my.ResBody)
		my.ReqHash1 = utils.HashString(my.requestText())

		result = append(result, my)
	}

	return result, nil
}
