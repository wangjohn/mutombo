package main

import (
	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/requests/:request_id", GetRequest)
	goji.Post("/requests", PostRequest)
	goji.Serve()
}
