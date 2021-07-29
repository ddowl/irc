package api

import (
	"encoding/json"
	"fmt"
	"io"
	"irc/server/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TODO: refactor validations in middleware fns
func SetRoutes() {
	http.HandleFunc("/api/rooms", ChatRoomsHandler)
	http.HandleFunc("/api/rooms/*", ChatRoomHandler)
	http.HandleFunc("/api/rooms/*/members", MembersHandler)
	http.HandleFunc("/api/rooms/*/members/*", MemberHandler)
	http.HandleFunc("/api/rooms/*/members/*/messages", MessagesHandler)
}

func ChatRoomsHandler(w http.ResponseWriter, r *http.Request) {
	store := model.GetChatRoomStore()

	switch r.Method {
	case http.MethodGet:
		ListChatRooms(w, store)
	case http.MethodPost:
		CreateChatRoom(w, store, r.Body)
	default:
		notFound(w, "Invalid HTTP Header")
	}
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomId, err := getRoomId(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	store := model.GetChatRoomStore()
	_, err = getChatRoom(roomId)
	if err != nil {
		badRequest(w, err)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		DeleteChatRoom(w, store, roomId)
	default:
		notFound(w, "Invalid HTTP Method")
	}
}

func MembersHandler(w http.ResponseWriter, r *http.Request) {
	roomId, err := getRoomId(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	room, err := getChatRoom(roomId)
	if err != nil {
		badRequest(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		notFound(w, "List chat room members not implemented")
	case http.MethodPost:
		JoinChatRoom(w, room, r.Body)
	default:
		notFound(w, "Invalid HTTP Method")
	}

	w.WriteHeader(200)
}

func MemberHandler(w http.ResponseWriter, r *http.Request) {
	roomId, err := getRoomId(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	tag, err := getMemberTag(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	room, err := getChatRoom(roomId)
	if err != nil {
		badRequest(w, err)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		LeaveChatRoom(w, room, tag)
	default:
		notFound(w, "Invalid HTTP Method")
	}
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	roomId, err := getRoomId(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	tag, err := getMemberTag(r.URL)
	if err != nil {
		badRequest(w, err)
		return
	}

	room, err := getChatRoom(roomId)
	if err != nil {
		badRequest(w, err)
		return
	}

	if !room.HasJoined(tag) {
		badRequest(w, fmt.Errorf(`"%s" hasn't joined room %+v`, tag, *room.GetMetadata()))
	}

	switch r.Method {
	case http.MethodPost:
		PostMessage(w, room, tag, r.Body)
	default:
		notFound(w, "Invalid HTTP Method")
	}

	w.WriteHeader(200)
}

func ListChatRooms(w http.ResponseWriter, store model.MessageProxyStore) {
	metadata := store.GetMetadata()
	res, err := json.Marshal(metadata)
	if err != nil {
		unexpectedError(w, err)
		return
	}

	w.Write(res)
}

type CreateChatRoomArgs struct {
	Name string `json:"name"`
}

type CreateChatRoomResponseBody struct {
	RoomId int `json:"roomId"`
}

func CreateChatRoom(w http.ResponseWriter, store model.MessageProxyStore, reqBody io.ReadCloser) {
	var args CreateChatRoomArgs
	err := json.NewDecoder(reqBody).Decode(&args)
	if err != nil {
		badRequest(w, err)
		return
	}

	roomId, err := store.AddProxy(args.Name)
	if err != nil {
		badRequest(w, err)
		return
	}

	body, err := json.Marshal(CreateChatRoomResponseBody{roomId})
	if err != nil {
		unexpectedError(w, err)
		return
	}

	w.Write(body)
}

type JoinChatRoomArgs struct {
	Tag         string `json:"tag"`
	CallbackURL string `json:"callbackUrl"`
}

func JoinChatRoom(w http.ResponseWriter, proxy model.MessageProxy, reqBody io.ReadCloser) {
	var args JoinChatRoomArgs
	err := json.NewDecoder(reqBody).Decode(&args)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = proxy.Join(args.Tag, args.CallbackURL)
	if err != nil {
		badRequest(w, err)
		return
	}
}

func LeaveChatRoom(w http.ResponseWriter, proxy model.MessageProxy, tag string) {
	err := proxy.Leave(tag)
	if err != nil {
		unexpectedError(w, err)
		return
	}
}

type PostMessageArgs struct {
	Message string `json:"message"`
}

func PostMessage(w http.ResponseWriter, proxy model.MessageProxy, memberTag string, reqBody io.ReadCloser) {
	var args PostMessageArgs
	err := json.NewDecoder(reqBody).Decode(&args)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = proxy.PostMessage(memberTag, args.Message)
	if err != nil {
		unexpectedError(w, err)
		return
	}
}

func DeleteChatRoom(w http.ResponseWriter, store model.MessageProxyStore, id int) {
	err := store.DeleteProxy(id)
	if err != nil {
		unexpectedError(w, err)
		return
	}
}

// Helpers / Validation

func getRoomId(url *url.URL) (int, error) {
	urlPathParams := strings.Split(url.Path, "/")
	if len(urlPathParams) < 4 {
		return 0, fmt.Errorf("url path doesn't contain room ID")
	}
	roomIdStr := urlPathParams[3]
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		return 0, fmt.Errorf(`chat room ID must be an integer: "%s"`, roomIdStr)
	}
	return roomId, nil
}

func getMemberTag(url *url.URL) (string, error) {
	urlPathParams := strings.Split(url.Path, "/")
	if len(urlPathParams) < 6 {
		return "", fmt.Errorf("url path doesn't contain room ID")
	}
	return urlPathParams[5], nil
}

func getChatRoom(id int) (model.MessageProxy, error) {
	store := model.GetChatRoomStore()
	room, err := store.GetProxy(id)
	if err != nil {
		return nil, fmt.Errorf(`chat room does not exist: "%d"`, id)
	}
	return room, nil
}

// Status Code Helpers

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	fmt.Fprint(w, err.Error())
}

func notFound(w http.ResponseWriter, err string) {
	w.WriteHeader(404)
	fmt.Fprint(w, err)
}

func unexpectedError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, err.Error())
}
