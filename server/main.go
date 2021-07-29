package main

import (
	"irc/server/api"
	"irc/server/model"
	"log"
	"net/http"
)

func main() {
	model.InitChatRoomStore()
	api.SetRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
