package requests

import (
	"encoding/json"
	"fmt"
	"linn221/Requester/utils"
	"net/url"
	"strings"
)

type MyRequest struct {
	Sequence int
	URL      string
	Method   string
	Domain   string

	ReqHeaders HeaderSlice
	ReqBody    string

	ResStatus  int
	ResHeaders HeaderSlice
	ResBody    string
	RespSize   int
	LatencyMs  int64

	RequestTime string
	// hashes
	ReqHash1    string // hash raw request
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

type Header struct {
	Name  string
	Value string
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s\n", h.Name, h.Value)
}

func (r MyRequest) requestText() string {

	raw := r.Method + " " + r.URL + " " + r.ReqBody + " " + r.ReqHeaders.EchoAll()
	return raw
}

func ParseHAR(bs []byte, resHashFunc func(*MyRequest) (string, string)) ([]MyRequest, error) {

	var har HAR
	if err := json.Unmarshal(bs, &har); err != nil {
		return nil, err
	}

	var result []MyRequest

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

		my := MyRequest{
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
