package requests

import (
	"encoding/json"
	"fmt"
	"linn221/Requester/utils"
	"net/url"
)

type MyRequest struct {
	Sequence int
	URL      string
	Method   string
	Domain   string

	ReqHeaders []Header
	ReqBody    string

	ResStatus  int
	ResHeaders []Header
	ResBody    string
	RespSize   int
	LatencyMs  int64

	ReqHash      string
	ResBodyHash  string
	ResTotalHash string
	RequestTime  string
}

type Header struct {
	Name  string
	Value string
}

func ParseHAR(bs []byte, hashFunc func(*MyRequest) string) ([]MyRequest, error) {

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

		requestHash := utils.HashString(fmt.Sprintf("%s %s %s", my.URL, my.Method, my.ReqBody))
		responseBodyHash := utils.HashString(my.ResBody)
		responseTotalHash := utils.HashString(hashFunc(&my))

		my.ReqHash = requestHash
		my.ResBodyHash = responseBodyHash
		my.ResTotalHash = responseTotalHash
		result = append(result, my)
	}

	return result, nil
}
