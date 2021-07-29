package model

import "fmt"

type ChatRoomStore struct {
	roomCounter int
	chatRooms   map[int]*ChatRoom
}

var state MessageProxyStore

func NewChatRoomStore() *ChatRoomStore {
	return &ChatRoomStore{roomCounter: 0, chatRooms: make(map[int]*ChatRoom)}
}

func InitChatRoomStore() {
	state = NewChatRoomStore()
}

func GetChatRoomStore() MessageProxyStore {
	return state
}

func (s *ChatRoomStore) AddProxy(name string) (int, error) {
	if !s.hasUniqueChatRoomName(name) {
		return 0, fmt.Errorf("cannot create duplicate chat room: %q", name)
	}

	id := s.roomCounter
	s.roomCounter += 1
	s.chatRooms[id] = EmptyChatRoom(name)
	return id, nil
}

func (s *ChatRoomStore) GetProxyMetadata() map[int]string {
	rooms := make(map[int]string)
	for roomId, room := range s.chatRooms {
		rooms[roomId] = room.name
	}
	return rooms
}

func (s *ChatRoomStore) GetProxy(roomId int) (MessageProxy, error) {
	if room, ok := s.chatRooms[roomId]; !ok {
		return nil, fmt.Errorf("chat room does not exist: %d", roomId)
	} else {
		return room, nil
	}
}

func (s *ChatRoomStore) DeleteProxy(roomId int) error {
	if _, ok := s.chatRooms[roomId]; !ok {
		return fmt.Errorf("chat room does not exist: %d", roomId)
	} else {
		// TODO: cleanup/flush chat room resources?
		delete(s.chatRooms, roomId)
		return nil
	}
}

func (s *ChatRoomStore) hasUniqueChatRoomName(name string) bool {
	for _, room := range s.chatRooms {
		if room.name == name {
			return false
		}
	}
	return true
}
