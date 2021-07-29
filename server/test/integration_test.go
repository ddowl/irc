package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"irc/server/api"
	"irc/server/model"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)

func listRoomsRequest() *http.Request {
	return httptest.NewRequest("GET", "/api/rooms", nil)
}

func createRoomRequest(roomName string) *http.Request {
	bs, err := json.Marshal(api.CreateChatRoomArgs{Name: roomName})
	if err != nil {
		log.Panicln(err)
	}
	return httptest.NewRequest("POST", "/api/rooms", bytes.NewReader(bs))
}

func joinRoomRequest(roomId int, tag string, callbackUrl string) *http.Request {
	bs, err := json.Marshal(api.JoinChatRoomArgs{tag, callbackUrl})
	if err != nil {
		log.Panicln(err)
	}
	return httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%d/members", roomId), bytes.NewReader(bs))
}

func leaveRoomRequest(roomId int, tag string) *http.Request {
	return httptest.NewRequest("DELETE", fmt.Sprintf("/api/rooms/%d/members/%s", roomId, tag), nil)
}

func postMessageRequest(roomId int, tag string, message string) *http.Request {
	bs, err := json.Marshal(api.PostMessageArgs{message})
	if err != nil {
		log.Panicln(err)
	}
	return httptest.NewRequest("POST", fmt.Sprintf("/api/rooms/%d/members/%s/messages", roomId, tag), bytes.NewReader(bs))
}

func deleteRoomRequest(roomId int) *http.Request {
	return httptest.NewRequest("DELETE", fmt.Sprintf(`/api/rooms/%d`, roomId), nil)
}

func TestListChatRoomsHandler(t *testing.T) {
	t.Run("new server has no chat rooms", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, listRoomsRequest())

		expectStatus(t, rr, 200)
		expectBody(t, rr, "[]")
	})

	t.Run("server with several rooms", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest("room0"))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.ChatRoomsHandler, createRoomRequest("room1"))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.ChatRoomsHandler, createRoomRequest("room2"))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.ChatRoomsHandler, listRoomsRequest())

		expectStatus(t, rr, 200)
		expectBody(t, rr, `[{"id":0,"name":"room0"},{"id":1,"name":"room1"},{"id":2,"name":"room2"}]`)
	})
}

func TestCreateChatRoomHandler(t *testing.T) {
	roomName := "room0"

	t.Run("create new room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))

		expectStatus(t, rr, 200)
		expectBody(t, rr, `{"roomId":0}`)
	})

	t.Run("enforces unique room names", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`cannot create duplicate chat room: "%s"`, roomName))
	})
}

func TestJoinChatRoomHandler(t *testing.T) {
	roomId := 0
	roomName := "room0"
	userTag := "new_user"
	callbackUrl := "localhost:6000"

	t.Run("join non-existent room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, callbackUrl))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`chat room does not exist: "%d"`, roomId))
	})

	t.Run("join empty room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, callbackUrl))

		expectStatus(t, rr, 200)
		expectBody(t, rr, "")
	})

	t.Run("join twice", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, callbackUrl))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, callbackUrl))
		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`"%s" already joined chat room {Id:%d Name:%s}`, userTag, roomId, roomName))
	})
}

func TestLeaveChatRoomHandler(t *testing.T) {
	roomId := 0
	roomName := "room0"
	userTag := "new_user"
	callbackUrl := "localhost:6000"

	t.Run("leave non-existent room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.MemberHandler, leaveRoomRequest(roomId, userTag))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`chat room does not exist: "%d"`, roomId))
	})

	t.Run("leave existing room, but haven't joined", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		leaveRoomId := 3
		rr = invokeHandler(api.MemberHandler, leaveRoomRequest(leaveRoomId, userTag))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`chat room does not exist: "%d"`, leaveRoomId))
	})

	t.Run("leave joined room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, callbackUrl))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MemberHandler, leaveRoomRequest(roomId, userTag))

		expectStatus(t, rr, 200)
		expectBody(t, rr, "")
	})
}

func TestPostMessageHandler(t *testing.T) {
	roomId := 0
	roomName := "room1"
	userTag := "new_user"
	message := "this is a test message"

	t.Run("post message to non-existent room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.MessagesHandler, postMessageRequest(roomId, userTag, message))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`chat room does not exist: "%d"`, roomId))
	})

	t.Run("post message to existing room, no one joined", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MessagesHandler, postMessageRequest(roomId, userTag, message))
		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`"%s" hasn't joined room {Id:%d Name:%s}`, userTag, roomId, roomName))
	})

	t.Run("post message to existing room, only you joined", func(t *testing.T) {
		model.InitChatRoomStore()

		ts := testServerExpectsNoCall(t)
		defer ts.Close()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, ts.URL))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MessagesHandler, postMessageRequest(roomId, userTag, message))
		expectStatus(t, rr, 200)
		expectBody(t, rr, "")
	})

	t.Run("post message to existing room, several people joined", func(t *testing.T) {
		model.InitChatRoomStore()

		myServer := testServerExpectsNoCall(t)
		defer myServer.Close()

		var wg sync.WaitGroup
		wg.Add(3)
		otherServer1 := testServerExpectsCall(t, &wg, message)
		defer otherServer1.Close()
		otherServer2 := testServerExpectsCall(t, &wg, message)
		defer otherServer2.Close()
		otherServer3 := testServerExpectsCall(t, &wg, message)
		defer otherServer3.Close()
		otherChatRoomMembers := []struct {
			tag         string
			callbackUrl string
		}{
			{tag: "user1", callbackUrl: otherServer1.URL},
			{tag: "user2", callbackUrl: otherServer1.URL},
			{tag: "user3", callbackUrl: otherServer1.URL},
		}

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, userTag, myServer.URL))
		expectStatus(t, rr, 200)

		for _, otherMember := range otherChatRoomMembers {
			rr = invokeHandler(api.MembersHandler, joinRoomRequest(roomId, otherMember.tag, otherMember.callbackUrl))
			expectStatus(t, rr, 200)
		}

		rr = invokeHandler(api.MessagesHandler, postMessageRequest(roomId, userTag, message))
		expectStatus(t, rr, 200)
		expectBody(t, rr, "")

		// If this test terminates, that means all other members' servers were called
		wg.Wait()
	})
}

func TestDeleteChatRoomHandler(t *testing.T) {
	roomId := 0
	roomName := "room0"

	t.Run("delete non-existent room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomHandler, deleteRoomRequest(roomId))

		expectStatus(t, rr, 400)
		expectBody(t, rr, fmt.Sprintf(`chat room does not exist: "%d"`, roomId))
	})

	t.Run("delete existing room", func(t *testing.T) {
		model.InitChatRoomStore()

		rr := invokeHandler(api.ChatRoomsHandler, createRoomRequest(roomName))
		expectStatus(t, rr, 200)

		rr = invokeHandler(api.ChatRoomHandler, deleteRoomRequest(roomId))

		expectStatus(t, rr, 200)
		expectBody(t, rr, "")
	})
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
	t.Helper()
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func expectBody(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	// Check the response body is what we expect.
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body:\n\tgot:  %v\n\twant: %v",
			rr.Body.String(), expected)
	}
}

func testServerExpectsNoCall(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "expected no callback to message poster")
	}))
}

func testServerExpectsCall(t *testing.T, wg *sync.WaitGroup, expectedMessage string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("called server!")
		assert.Equal(t, r.Method, http.MethodPost)

		var body model.CallbackBody
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.Nil(t, err)
		assert.Equal(t, expectedMessage, body.Message)
		wg.Done()
	}))
}
