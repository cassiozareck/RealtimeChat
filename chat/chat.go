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
	people   []shared.Person
}

// GetChat returns a Chat object with the given id. If the chat does not exist,
// it returns an error. If the chat exists, it returns a Chat object with the
// messages loaded.
func GetChat(chatDB db.ChatDB, ID uint32) (*Chat, error) {
	c := Chat{ID, chatDB, []shared.Message{}, []shared.Person{}}

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
	err = c.chatDB.Store(c.id, message)

	if err != nil {
		return err
	}
	return nil
}

func (c *Chat) GetMessages() []shared.Message {
	return c.messages
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
// stored in the database.
func (c *Chat) UpdatePeople() error {
	messages, err := c.chatDB.GetMessages(c.id)

	if err != nil {
		return err
	}

	for _, msg := range messages {
		// Create a logic that checks if the person is already in the list
		// if not, add it to the list
		if !c.checkPerson(msg.UserID) {
			c.people = append(c.people, shared.Person{ID: msg.UserID})
		}
	}
	return nil
}

func (c *Chat) GetPeople() []shared.Person {
	return c.people
}

func (c *Chat) LastMessage() (*shared.Message, error) {
	messages, err := c.chatDB.GetMessages(c.id)
	if err != nil {
		return nil, err
	}
	return &messages[len(messages)-1], nil
}

func (c *Chat) GetID() uint32 {
	return c.id
}

func (c *Chat) checkPerson(ID uint32) bool {
	for _, person := range c.people {
		if person.ID == ID {
			return true
		}
	}
	return false
}

func (c *Chat) checkSenderID(SenderID uint32) error {
	if SenderID <= 0 {
		return fmt.Errorf("invalid sender id: %d", SenderID)
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
