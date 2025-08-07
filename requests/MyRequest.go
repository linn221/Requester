package requests

import (
	"encoding/json"
	"io/ioutil"
	"linn221/Requester/har"
	"net/url"
)

type MyRequest struct {
	Seq        int
	URL        string
	Domain     string
	ReqHeaders []Header
	ReqBody    string
	ResHeaders []Header
	ResBody    string
	LatencyMs  int64
}

type Header struct {
	Name  string
	Value string
}

func ParseHAR(path string) ([]MyRequest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var har har.HAR
	if err := json.Unmarshal(data, &har); err != nil {
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
			Seq:        i + 1,
			URL:        entry.Request.URL,
			Domain:     domain,
			ReqHeaders: reqHeaders,
			ReqBody:    entry.Request.PostData.Text,
			ResHeaders: resHeaders,
			ResBody:    resBody,
			LatencyMs:  int64(entry.Time),
		}
		result = append(result, my)
	}

	return result, nil
}
