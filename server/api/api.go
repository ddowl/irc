package api

import (
	"encoding/json"
	"fmt"
	"irc/server/model"
	"net/http"
)

func SetRoutes() {
	http.HandleFunc("/api/rooms", ChatRoomsHandler)
	http.HandleFunc("/api/rooms/*/members", MembersHandler)
	http.HandleFunc("/api/rooms/*/messages", MessagesHandler)
}

func ChatRoomsHandler(w http.ResponseWriter, r *http.Request) {
	state := model.GetChatRoomStore()

	switch r.Method {
	case http.MethodGet:
		ListChatRooms(w, r, state)
	case http.MethodPost:
		CreateChatRoom(w, r, state)
	case http.MethodDelete:
		DeleteChatRoom(w, r, state)
	default:
		w.WriteHeader(200)
		fmt.Fprintln(w, "Not found")
	}
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func ListChatRooms(w http.ResponseWriter, r *http.Request, m model.MessageProxyStore) {
	res, err := json.Marshal(m.GetProxyMetadata())
	if err != nil {
		unexpectedError(w, err)
		return
	}

	w.Write(res)
	w.WriteHeader(200)
}

type CreateChatRoomArgs struct {
	Name string `json:"name"`
}

type CreateChatRoomResponseBody struct {
	RoomId int `json:"roomId"`
}

func CreateChatRoom(w http.ResponseWriter, r *http.Request, m model.MessageProxyStore) {
	var args CreateChatRoomArgs
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		unexpectedError(w, err)
		return
	}

	roomId, err := m.AddProxy(args.Name)
	if err != nil {
		badRequest(w, err)
		return
	}

	w.WriteHeader(200)
	body, err := json.Marshal(CreateChatRoomResponseBody{roomId})
	if err != nil {
		unexpectedError(w, err)
		return
	}

	w.Write(body)
}

func JoinChatRoom(w http.ResponseWriter, s *model.ChatRoomStore) {
	w.WriteHeader(200)
}

func LeaveChatRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func DeleteChatRoom(w http.ResponseWriter, r *http.Request, m model.MessageProxyStore) {
	w.WriteHeader(200)
}

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	fmt.Fprintln(w, err.Error())
}

func unexpectedError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintln(w, err.Error())
}
