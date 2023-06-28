package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
)

type Chat struct {
	chatDB   db.ChatDB
	userID   uint32
	senderID uint32
}

func GetChat(chatDB db.ChatDB, userID, destinID uint32) *Chat {

	c := Chat{chatDB, userID, destinID}
	return &c
}

func (c *Chat) SendMessage(msg string) error {
	if msg == "" {
		return fmt.Errorf("message cannot be empty")
	}
	chatExist, err := c.chatDB.ChatExists(c.userID)

	if err != nil {
		return err
	}

	if !chatExist {
		return fmt.Errorf("chat does not exist")
	}

	message := shared.NewMessage(c.userID, msg)
	err = c.chatDB.Store(message)

	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) CreateChat() (uint32, error) {
	id, err := c.chatDB.CreateChat()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Chat) GetMessages() ([]shared.Message, error) {
	msgs, err := c.chatDB.GetMessages()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Chat) LastMessage() (*shared.Message, error) {
	msgs, err := c.chatDB.GetMessages()
	if err != nil {
		return nil, err
	}
	return &msgs[len(msgs)-1], nil
}
