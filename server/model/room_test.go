package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const roomId = 0
const roomName = "test_chat_room"
const userName = "new_user"
const callbackUrl = "https://dummy.com:8080"

func TestJoin(t *testing.T) {
	t.Run("join empty room", func(t *testing.T) {
		room := EmptyChatRoom(roomId, roomName)

		err := room.Join(userName, callbackUrl)

		assert.Nil(t, err)
		assert.Equal(t, callbackUrl, room.members[userName])
	})

	t.Run("join twice", func(t *testing.T) {
		room := EmptyChatRoom(roomId, roomName)

		err := room.Join(userName, callbackUrl)
		assert.Nil(t, err)

		err = room.Join(userName, callbackUrl)
		assert.NotNil(t, err)

		assert.Equal(t, callbackUrl, room.members[userName])
	})
}

func TestLeave(t *testing.T) {
	t.Run("leave empty room", func(t *testing.T) {
		room := EmptyChatRoom(roomId, roomName)

		err := room.Leave(userName)

		assert.NotNil(t, err)

		_, ok := room.members[userName]
		assert.False(t, ok)
	})

	t.Run("leave joined room", func(t *testing.T) {
		room := EmptyChatRoom(roomId, roomName)

		err := room.Join(userName, callbackUrl)
		assert.Nil(t, err)

		err = room.Leave(userName)
		assert.Nil(t, err)

		_, ok := room.members[userName]
		assert.False(t, ok)
	})
}
