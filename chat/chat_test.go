package chat

import (
	"github.com/cassiozareck/realchat/shared"
	"testing"
)

type ChatDBMock struct {
	messages []shared.Message
}

func (c *ChatDBMock) CreateChat() (uint32, error) {
	return 0, nil
}

func (c *ChatDBMock) ChatExists(chatID uint32) (bool, error) {
	return true, nil
}

func (c *ChatDBMock) Store(msg shared.Message) error {
	c.messages = append(c.messages, msg)
	return nil
}

func (c *ChatDBMock) GetMessages() ([]shared.Message, error) {
	return c.messages, nil
}

// TestChat_SendMSG will send a message and test if the
// last message is giving the appropriate answer
func TestChat_SendMSG(t *testing.T) {

	chatDB := ChatDBMock{}

	c := GetChat(&chatDB, 123, 321)

	MSG := "hello"
	err := c.SendMessage(MSG)
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
	destinID := uint32(123)
	chatDB := ChatDBMock{messages: []shared.Message{
		shared.NewMessage(destinID, "hello"),
		shared.NewMessage(destinID, "how are u?"),
		shared.NewMessage(userID, "im fine thanks"),
	}}

	c := GetChat(&chatDB, userID, destinID)
	messages, err := c.GetMessages()
	if err != nil {
		t.Fatal("Could not get messages")
	}
	if messages[0].Text != "hello" ||
		messages[1].Text != "how are u?" ||
		messages[2].Text != "im fine thanks" {
		t.Fatal("Messages retrieved not the same as database")
	}
	if messages[0].UserID != destinID ||
		messages[2].UserID != userID {
		t.Fatal("Messages retrieved not the same ID as database")
	}

}
