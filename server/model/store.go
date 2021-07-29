package model

import (
	"fmt"
	"sort"
)

type ChatRoomStore struct {
	roomCounter int
	chatRooms   map[int]*ChatRoom
}

var store *ChatRoomStore

func NewChatRoomStore() *ChatRoomStore {
	return &ChatRoomStore{roomCounter: 0, chatRooms: make(map[int]*ChatRoom)}
}

func InitChatRoomStore() {
	store = NewChatRoomStore()
}

func GetChatRoomStore() *ChatRoomStore {
	return store
}

func (s *ChatRoomStore) AddProxy(name string) (int, error) {
	if !s.hasUniqueChatRoomName(name) {
		return 0, fmt.Errorf("cannot create duplicate chat room: %q", name)
	}

	id := s.roomCounter
	s.roomCounter += 1
	s.chatRooms[id] = EmptyChatRoom(id, name)
	return id, nil
}

func (s *ChatRoomStore) GetMetadata() []ProxyMetadata {
	// Ensures that chat room metadata is retrieved sorted by room ID
	keys := make([]int, 0, len(s.chatRooms))
	for k := range s.chatRooms {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	rooms := make([]ProxyMetadata, 0, len(s.chatRooms))
	for k := range keys {
		room := s.chatRooms[k]
		rooms = append(rooms, *room.GetMetadata())
	}
	return rooms
}

func (s *ChatRoomStore) GetProxy(id int) (MessageProxy, error) {
	if room, ok := s.chatRooms[id]; !ok {
		return nil, fmt.Errorf("chat room does not exist: %d", id)
	} else {
		return room, nil
	}
}

func (s *ChatRoomStore) DeleteProxy(id int) error {
	if _, ok := s.chatRooms[id]; !ok {
		return fmt.Errorf("chat room does not exist: %d", id)
	} else {
		// TODO: cleanup/flush chat room resources?
		delete(s.chatRooms, id)
		return nil
	}
}

func (s *ChatRoomStore) hasUniqueChatRoomName(name string) bool {
	for _, room := range s.chatRooms {
		if room.Name == name {
			return false
		}
	}
	return true
}
