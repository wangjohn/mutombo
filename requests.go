package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wangjohn/gowebutils"
	"github.com/zenazn/goji/web"
)

type RequestData struct {
	Body     interface{} `json:"body"`
	Headers  interface{} `json:"headers"`
	URL      string      `json:"url"`
	Method   string      `json:"method"`
	Blocking bool        `json:"blocking"`
}

func GetRequest(c web.C, w http.ResponseWriter, r *http.Request) {
	requestId := c.URlParams["request_id"]
	if requestId == "" {
		err = errors.New("You must specify a request_id when retrieving a request")
		gowebutils.SendError(w, err)
		return
	}
}

const (
	maxClientTimeoutMinutes = 30
)

func PostRequest(w http.ResponseWriter, r *http.Request) {
	store, body, err := prepareRequest(r)
	defer store.Close()
	if err != nil {
		gowebutils.SendError(w, err)
		return
	}

	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		gowebutils.SendError(w, err)
		return
	}

	resp, err := makeRequest(data)
	for h, _ := range resp.Header {
		w.Header().Set(h, resp.Header.Get(h))
	}

	w.Write(resp.Body)
}

// Makes the actual HTTP request and returns a response and/or error.
func makeRequest(data RequestData) (*http.Response, error) {
	var method string
	switch data.Method {
	case "GET", "POST":
		method = data.Method
	default:
		return nil, fmt.Errorf("Invalid method type '%v'", data.Method)
	}

	body, err := json.Marshal(data.Body)
	if err != nil {
		return nil, err
	}
	httpReq := http.NewRequest(data.Method, data.URL, body)

	for headerName, headerVal := range data.Headers {
		switch hv := headerVal.(type) {
		case string:
			httpReq.Header.Add(headerName, hv)
		default:
			return nil, fmt.Errorf("Invalid header value for '%v'. Header values must be strings", headerName)
		}
	}

	client := &http.Client{
		Timeout: time.Duration(maxClientTimeoutMinutes) * time.Minute,
	}
	return client.Do(httpReq)
}
