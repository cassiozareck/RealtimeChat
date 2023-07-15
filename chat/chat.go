package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
)

type Chat struct {
	id     uint32
	chatDB db.ChatDB
}

// GetChat returns a Chat object with the given id. If the chat does not exist,
// it returns an error. If the chat exists, it returns a Chat object with the
// messages loaded.
func GetChat(chatDB db.ChatDB, ID uint32) (*Chat, error) {
	exist, err := chatDB.ChatExists(ID)

	if err != nil {
		return nil, err
	}

	if exist {
		c := Chat{chatDB: chatDB, id: ID}
		return &c, nil
	}

	return nil, fmt.Errorf("chat with id %d does not exist", ID)
}

// NewChat creates a new chat with unique id and returns a Chat object.
func NewChat(chatDB db.ChatDB) (*Chat, error) {
	id, err := chatDB.CreateChat()
	if err != nil {
		return nil, err
	}
	c := Chat{chatDB: chatDB, id: id}

	return &c, nil
}

func (c *Chat) SendMessage(message shared.IncomingMessage) error {

	newMessage, err := shared.NewMessage(
		message.SenderID,
		message.ChatID,
		message.Text,
	)
	if err != nil {
		return err
	}

	err = c.chatDB.Store(newMessage)

	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) GetMessages() ([]shared.Message, error) {
	messages, err := c.chatDB.GetMessages(c.id)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetPeople returns a list of people ids that are in the chat.
func (c *Chat) GetPeople() ([]uint32, error) {
	people := make([]uint32, 0)

	messages, err := c.GetMessages()
	if err != nil {
		return nil, err
	}

	for _, msg := range messages {
		// Create a logic that checks if the person is already in the list
		// if not, add it to the list
		if !shared.Contains(people, msg.SenderID) {
			people = append(people, msg.SenderID)
		}
	}

	return people, nil
}

func (c *Chat) GetID() uint32 {
	return c.id
}
