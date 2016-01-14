package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wangjohn/gowebutils"
	"github.com/wangjohn/mutombo/storage"
	"github.com/zenazn/goji/web"
)

type RequestData struct {
	Body     interface{}       `json:"body"`
	Headers  map[string]string `json:"headers"`
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Blocking bool              `json:"blocking"`
}

type NonBlockingPostResp struct {
	RequestId string `json:"request_id"`
}

func GetRequest(c web.C, w http.ResponseWriter, r *http.Request) {
	requestId := c.URLParams["request_id"]
	if requestId == "" {
		err := errors.New("You must specify a request_id when retrieving a request")
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

	if data.Blocking {
		// Blocking request
		gowebutils.SendError(w, fmt.Errorf("Blocking requests are currently unsupported"))
		return
	} else {
		// Non-blocking request
		storedRequest, err := store.StoreRequest(data.Blocking, data.Method, data.URL)
		if err != nil {
			gowebutils.SendError(w, err)
			return
		}
		if storedRequest == nil {
			gowebutils.SendError(w, fmt.Errorf("Unable to store non-blocking request"))
			return
		}
		err = respondToPost(w, storedRequest)
		if err != nil {
			gowebutils.SendError(w, err)
			return
		}

		resp, err := makeRequest(data)
		_, err = store.StoreResponse(storedRequest.RequestId, resp)
		if err != nil {
			log.Printf("Unable to store response: %v", err)
		}
	}
}

func respondToPost(w http.ResponseWriter, storedRequest *storage.StoredRequest) error {
	resp := NonBlockingPostResp{
		RequestId: storedRequest.RequestId,
	}
	return json.NewEncoder(w).Encode(resp)
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
	bodyReader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(method, data.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for name, val := range data.Headers {
		httpReq.Header.Add(name, val)
	}

	client := &http.Client{
		Timeout: time.Duration(maxClientTimeoutMinutes) * time.Minute,
	}
	return client.Do(httpReq)
}
