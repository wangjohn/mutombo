package main

import (
	"log"
	"net/http"

	"github.com/wangjohn/gowebutils"
	"github.com/wangjohn/mutombo/storage"
	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/requests/:request_id", GetRequest)
	goji.Post("/requests", PostRequest)
	goji.Serve()
}

func prepareRequest(r *http.Request) (storage.Storage, []byte, error) {
	store, err := storage.GenerateStorage(storage.Postgres)
	if err != nil {
		return store, nil, err
	}
	body, err := gowebutils.PrepareRequestBody(r)
	log.Printf("Request body: %v", string(body))
	return store, body, err
}
