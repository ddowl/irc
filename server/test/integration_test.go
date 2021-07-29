package test

import (
	"irc/server/api"
	"irc/server/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func listRoomsRequest() *http.Request {
	return httptest.NewRequest("GET", "/api/rooms", nil)
}

func createRoomRequest(roomName string) *http.Request {
	return httptest.NewRequest("POST", "/api/rooms", strings.NewReader(roomName))
}

func TestListChatroomsHandler(t *testing.T) {
	t.Run("new server has no chat rooms", func(t *testing.T) {
		model.InitChatRoomStore()
		req := listRoomsRequest()

		rr := invokeHandler(api.ChatRoomsHandler, req)

		expectStatus(t, rr, 200)
		expectBody(t, rr, "{}")
	})

	//t.Run("server with several rooms", func(t *testing.T) {
	//	model.InitChatRoomStore()
	//	req := listRoomsRequest()
	//
	//	rr := invokeHandler(api.ChatRoomsHandler, req)
	//
	//	expectStatus(t, rr, 200)
	//	expectBody(t, rr, "{}")
	//})
}

func TestCreateChatroomHandler(t *testing.T) {
	model.InitChatRoomStore()
	req := httptest.NewRequest("POST", "/api/rooms", strings.NewReader(`{ "name": "room1" }`))

	rr := invokeHandler(api.ChatRoomsHandler, req)

	expectStatus(t, rr, 200)
	expectBody(t, rr, `{"roomId":0}`)
}

func TestJoinChatroomHandler(t *testing.T) {
	model.InitChatRoomStore()
	req := httptest.NewRequest("POST", "/api/rooms/1/members", nil)

	rr := invokeHandler(api.MembersHandler, req)

	expectStatus(t, rr, 200)
	expectBody(t, rr, "")
}

func TestLeaveChatroomHandler(t *testing.T) {
	model.InitChatRoomStore()
	req := httptest.NewRequest("DELETE", "/api/rooms/1/members", nil)

	rr := invokeHandler(api.MembersHandler, req)

	expectStatus(t, rr, 200)
	expectBody(t, rr, "")
}

func TestPostMessageHandler(t *testing.T) {
	model.InitChatRoomStore()
	req := httptest.NewRequest("POST", "/api/rooms/1/messages", nil)

	rr := invokeHandler(api.MessagesHandler, req)

	expectStatus(t, rr, 200)
	expectBody(t, rr, "")
}

func TestDeleteChatroomHandler(t *testing.T) {
	model.InitChatRoomStore()
	req := httptest.NewRequest("DELETE", "/api/rooms", nil)

	rr := invokeHandler(api.ChatRoomsHandler, req)

	expectStatus(t, rr, 200)
	expectBody(t, rr, "")
}

func invokeHandler(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	return rr
}

func expectStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func expectBody(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	// Check the response body is what we expect.
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
