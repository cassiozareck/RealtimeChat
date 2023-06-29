package chat

import (
	"github.com/cassiozareck/realchat/shared"
	"testing"
)

type ChatDBMock struct {
	ids      []uint32
	messages []shared.Message
}

func (c *ChatDBMock) CreateChat() (uint32, error) {
	c.ids = append(c.ids, c.ids[len(c.ids)-1]+1)
	return 0, nil
}

func (c *ChatDBMock) ChatExists(chatID uint32) (bool, error) {
	for _, id := range c.ids {
		if id == chatID {
			return true, nil
		}
	}
	return false, nil
}

func (c *ChatDBMock) Store(chatID uint32, msg shared.Message) error {
	c.ids = append(c.ids, chatID)
	c.messages = append(c.messages, msg)
	return nil
}

func (c *ChatDBMock) GetMessages(chatID uint32) ([]shared.Message, error) {
	exist, err := c.ChatExists(chatID)
	if exist {
		return c.messages, nil
	}
	return nil, err
}

const SenderId = 123

// TestChat_SendMSG will send a message and test if the
// last message is giving the appropriate answer
func TestChat_SendMSG(t *testing.T) {

	chatDB := ChatDBMock{}

	c, err := NewChat(&chatDB)
	if err != nil {
		t.Fatal("Error while creating new chat: ", err)
	}

	MSG := "hello"
	err = c.SendMessage(SenderId, MSG)
	if err != nil {
		t.Fatal("Error while sending message: ", err)
	}

	lastMessage, err := c.LastMessage()
	if err != nil {
		t.Fatal("Error while retrieving last message: ", err)
	}

	if lastMessage.Text != MSG {
		t.Fatal("Last Message not equals ", MSG)
	}
}

// TestChat_GetMessages will get messages using Chat and compare to see if
// messages are matching the same from database
func TestChat_GetMessages(t *testing.T) {
	userID := uint32(123)
	senderID := uint32(321)
	chatDB := ChatDBMock{messages: []shared.Message{
		shared.NewMessage(senderID, "hello"),
		shared.NewMessage(senderID, "how are u?"),
		shared.NewMessage(userID, "im fine thanks"),
	}}

	c, err := NewChat(&chatDB)
	if err != nil {
		t.Fatal("Error while creating new chat: ", err)
	}

	messages, err := c.GetMessages()
	if err != nil {
		t.Fatal("Could not get messages")
	}
	if messages[0].Text != "hello" ||
		messages[1].Text != "how are u?" ||
		messages[2].Text != "im fine thanks" {
		t.Fatal("Messages retrieved not the same as database")
	}
	if messages[0].UserID != senderID ||
		messages[2].UserID != userID {
		t.Fatal("Messages retrieved not the same ID as database")
	}
}

// TestChat_Exist will test if chat exists
func TestChat_Exist(t *testing.T) {
	chatDB := ChatDBMock{}

	_, err := GetChat(&chatDB, 123)
	if err == nil {
		t.Fatal("Chat should not exist")
	}
}
