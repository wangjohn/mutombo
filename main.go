package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/wangjohn/gowebutils"
	"github.com/wangjohn/mutombo/storage"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var postgresPassword = ""

func main() {
	flag.StringVar(&postgresPassword, "postgres_password", "", "The password for the postgres db")
	flag.Parse()
	goji.Get("/requests/:request_id", GetRequest)
	goji.Post("/requests", PostRequest)
	goji.Use(CORSHandler())
	goji.Serve()
}

func prepareRequest(r *http.Request) (storage.Storage, []byte, error) {
	store, err := storage.GenerateStorage(storage.Postgres, postgresPassword)
	if err != nil {
		return store, nil, err
	}
	body, err := gowebutils.PrepareRequestBody(r)
	log.Printf("Request body: %v", string(body))
	return store, body, err
}

func CORSHandler() web.MiddlewareType {
	var c *cors.Cors
	c = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})
	return c.Handler
}
