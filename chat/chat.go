package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
)

type Chat struct {
	ID     uint32
	chatDB db.ChatDB
}

// GetChat returns a Chat object with the given ID.
func GetChat(chatDB db.ChatDB, ID uint32) (*Chat, error) {
	c := Chat{ID, chatDB}
	exist, err := chatDB.ChatExists(ID)
	if err != nil {
		return nil, err
	}
	if exist {
		return &c, nil
	}
	return nil, fmt.Errorf("chat with ID %d does not exist", ID)
}

// NewChat creates a new chat with unique ID and returns a Chat object.
func NewChat(chatDB db.ChatDB) (*Chat, error) {
	c := Chat{0, chatDB}

	id, err := chatDB.CreateChat()
	if err != nil {
		return nil, err
	}
	c.ID = id

	return &c, nil
}

func (c *Chat) SendMessage(SenderID uint32, msg string) error {
	err := c.checkSenderID(SenderID)
	if err != nil {
		return err
	}
	err = c.checkMessage(msg)
	if err != nil {
		return err
	}

	message := shared.NewMessage(SenderID, msg)
	err = c.chatDB.Store(c.ID, message)

	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) GetMessages() ([]shared.Message, error) {
	msgs, err := c.chatDB.GetMessages(c.ID)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Chat) LastMessage() (*shared.Message, error) {
	msgs, err := c.chatDB.GetMessages(c.ID)
	if err != nil {
		return nil, err
	}
	return &msgs[len(msgs)-1], nil
}

func (c *Chat) GetID() uint32 {
	return c.ID
}

func (c *Chat) checkSenderID(SenderID uint32) error {
	if SenderID <= 0 {
		return fmt.Errorf("invalid sender ID: %d", SenderID)
	}
	return nil
}

func (c *Chat) checkMessage(msg string) error {
	if len(msg) == 0 {
		return fmt.Errorf("invalid message: %s", msg)
	}
	if len(msg) > shared.MAX_MESSAGE_SIZE {
		return fmt.Errorf("message too long: %s", msg)
	}
	return nil
}
