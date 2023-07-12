package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
)

type Chat struct {
	id       uint32
	chatDB   db.ChatDB
	messages []shared.Message
	people   []uint32
}

// GetChat returns a Chat object with the given id. If the chat does not exist,
// it returns an error. If the chat exists, it returns a Chat object with the
// messages loaded.
func GetChat(chatDB db.ChatDB, ID uint32) (*Chat, error) {
	c := Chat{chatDB: chatDB}

	exist, err := chatDB.ChatExists(ID)

	if err != nil {
		return nil, err
	}

	if exist {
		err := c.UpdateMessages()
		if err != nil {
			return nil, err
		}

		err = c.UpdatePeople()
		if err != nil {
			return nil, err
		}
		return &c, nil
	}

	return nil, fmt.Errorf("chat with id %d does not exist", ID)
}

// NewChat creates a new chat with unique id and returns a Chat object.
func NewChat(chatDB db.ChatDB) (*Chat, error) {
	c := Chat{chatDB: chatDB}

	id, err := chatDB.CreateChat()
	if err != nil {
		return nil, err
	}
	c.id = id

	return &c, nil
}

func (c *Chat) SendMessage(message shared.Message) error {

	err := c.chatDB.Store(message)

	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) GetMessages() []shared.Message {
	return c.messages
}

// GetPeople returns a list of people ids that are in the chat.
func (c *Chat) GetPeople() []uint32 {
	return c.people
}

func (c *Chat) GetID() uint32 {
	return c.id
}

func (c *Chat) UpdateMessages() error {
	messages, err := c.chatDB.GetMessages(c.id)
	if err != nil {
		return err
	}
	c.messages = messages
	return nil
}

// UpdatePeople updates the people list of the chat using the messages
// stored in the database. Important to note that it get people by using
// the messages, so if the messages are not updated, the people list will
// not be updated as well.
func (c *Chat) UpdatePeople() error {
	for _, msg := range c.messages {
		// Create a logic that checks if the person is already in the list
		// if not, add it to the list
		if !c.checkIfPersonExists(msg.SenderID()) {
			c.people = append(c.people, msg.SenderID())
		}
	}
	return nil
}

func (c *Chat) checkIfPersonExists(ID uint32) bool {
	for _, id := range c.people {
		if id == ID {
			return true
		}
	}
	return false
}
