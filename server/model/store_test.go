package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChatRoomMeta(t *testing.T) {
	t.Run("new server has no chat rooms", func(t *testing.T) {
		s := NewChatRoomStore()

		meta := s.GetProxyMetadata()

		assert.Equal(t, 0, len(meta))
	})

	t.Run("server with several rooms", func(t *testing.T) {
		rooms := make(map[int]*ChatRoom)
		rooms[0] = &ChatRoom{
			name: "room0",
		}
		rooms[1] = &ChatRoom{
			name: "room1",
		}
		rooms[2] = &ChatRoom{
			name: "room2",
		}

		s := ChatRoomStore{
			roomCounter: 3,
			chatRooms:   rooms,
		}

		meta := s.GetProxyMetadata()

		assert.Equal(t, 3, len(meta))
	})
}

func TestAddChatRoom(t *testing.T) {
	t.Run("new server will have 1 chat room", func(t *testing.T) {
		s := NewChatRoomStore()

		_, err := s.AddProxy("room1")

		assert.Nil(t, err)
		assert.Equal(t, len(s.chatRooms), 1)
	})

	t.Run("server with existing rooms will have 1 more", func(t *testing.T) {
		rooms := make(map[int]*ChatRoom)
		rooms[0] = &ChatRoom{
			name: "room0",
		}
		rooms[1] = &ChatRoom{
			name: "room1",
		}
		rooms[2] = &ChatRoom{
			name: "room2",
		}

		s := ChatRoomStore{
			roomCounter: 3,
			chatRooms:   rooms,
		}

		_, err := s.AddProxy("room3")

		assert.Nil(t, err)
		assert.Equal(t, len(s.chatRooms), 4)
	})

	t.Run("chat room names must be unique", func(t *testing.T) {
		s := NewChatRoomStore()

		_, err := s.AddProxy("room1")
		assert.Nil(t, err)

		_, err = s.AddProxy("room1")
		assert.NotNil(t, err)

		assert.Equal(t, len(s.chatRooms), 1)
	})
}

func TestGetChatRoom(t *testing.T) {
	t.Run("err if room ID not present", func(t *testing.T) {
		s := NewChatRoomStore()

		room, err := s.GetProxy(0)
		assert.NotNil(t, err)
		assert.Nil(t, room)
	})

	t.Run("returns room with associated ID", func(t *testing.T) {
		roomID := 5

		rooms := make(map[int]*ChatRoom)
		rooms[roomID] = &ChatRoom{
			name: "room5",
		}

		s := ChatRoomStore{
			chatRooms: rooms,
		}

		room, err := s.GetProxy(roomID)
		assert.Nil(t, err)
		assert.NotNil(t, room)
	})
}

func TestDeleteChatRoom(t *testing.T) {
	t.Run("err if room ID not present", func(t *testing.T) {
		s := NewChatRoomStore()

		err := s.DeleteProxy(0)
		assert.NotNil(t, err)
	})

	t.Run("deletes room with associated ID", func(t *testing.T) {
		roomID := 5

		rooms := make(map[int]*ChatRoom)
		rooms[roomID] = &ChatRoom{
			name: "room5",
		}

		s := ChatRoomStore{
			chatRooms: rooms,
		}

		err := s.DeleteProxy(roomID)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(s.chatRooms))
	})
}
