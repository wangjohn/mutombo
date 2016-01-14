package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wangjohn/gowebutils"
	"github.com/zenazn/goji/web"
)

type RequestData struct {
	Body     interface{}       `json:"body"`
	Headers  map[string]string `json:"headers"`
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Blocking bool              `json:"blocking"`
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

	resp, err := makeRequest(data)
	for h, _ := range resp.Header {
		w.Header().Set(h, resp.Header.Get(h))
	}

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gowebutils.SendError(w, err)
		return
	}
	w.Write(byteBody)
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
