package main

import (
	"irc/server/api"
	"log"
	"net/http"
)

func main() {
	api.SetRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
