package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatRoom struct {
	ProxyMetadata
	members map[string]string
}

type ChatRoomMetadata struct {
	id   int
	name string
}

func EmptyChatRoom(id int, name string) *ChatRoom {
	return &ChatRoom{ProxyMetadata{id, name}, make(map[string]string)}
}

func (c *ChatRoom) Join(tag string, callbackUrl string) error {
	if _, ok := c.members[tag]; ok {
		return fmt.Errorf(`"%s" already joined chat room %+v`, tag, *c.GetMetadata())
	}

	c.members[tag] = callbackUrl
	return nil
}

func (c *ChatRoom) HasJoined(tag string) bool {
	if _, ok := c.members[tag]; ok {
		return true
	}
	return false
}

func (c *ChatRoom) Leave(tag string) error {
	if _, ok := c.members[tag]; !ok {
		return fmt.Errorf(`"%s" is not in chat room "%s"`, tag, c.Name)
	}

	delete(c.members, tag)
	return nil
}

// TODO: move to controller layer?
type CallbackBody struct {
	Message string `json:"message"`
}

func (c *ChatRoom) PostMessage(tag string, message string) error {
	for member, callbackUrl := range c.members {
		if member != tag {
			// TODO: dispatch messages in parallel
			// TODO: log client responses and errors
			bs, err := json.Marshal(CallbackBody{Message: message})
			if err != nil {
				return err
			}
			http.Post(callbackUrl, "application/json", bytes.NewReader(bs))
		}
	}
	return nil
}

func (c *ChatRoom) GetMetadata() *ProxyMetadata {
	return &ProxyMetadata{c.Id, c.Name}
}
