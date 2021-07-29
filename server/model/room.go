package model

import (
	"fmt"
	"net/http"
	"strings"
)

type ChatRoom struct {
	name    string
	members map[string]string
}

func EmptyChatRoom(name string) *ChatRoom {
	return &ChatRoom{name, make(map[string]string)}
}

func (c *ChatRoom) Join(tag string, callbackUrl string) error {
	if _, ok := c.members[tag]; ok {
		return fmt.Errorf(`"%s" already joined chat room "%s"`, tag, c.name)
	}

	c.members[tag] = callbackUrl
	return nil

}

func (c *ChatRoom) Leave(tag string) error {
	if _, ok := c.members[tag]; !ok {
		return fmt.Errorf(`"%s" is not in chat room "%s"`, tag, c.name)
	}

	delete(c.members, tag)
	return nil
}

func (c *ChatRoom) PostMessage(tag string, message string) error {
	for member, callbackUrl := range c.members {
		if member != tag {
			// TODO: dispatch messages in parallel
			// TODO: log client responses and errors
			http.Post(callbackUrl, "text/plain", strings.NewReader(message))
		}
	}
	return nil
}
