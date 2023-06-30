package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/shared"
	"testing"
)

type ChatDBMock struct {
	chatId   uint32
	exist    bool
	messages []shared.Message
	err      error
}

func (c *ChatDBMock) CreateChat() (uint32, error) {
	return c.chatId, c.err
}

func (c *ChatDBMock) ChatExists(chatID uint32) (bool, error) {
	return c.exist, c.err
}

func (c *ChatDBMock) Store(chatID uint32, msg shared.Message) error {
	return c.err
}

func (c *ChatDBMock) GetMessages(chatID uint32) ([]shared.Message, error) {
	return c.messages, c.err
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
	chatID := 1

	testCases := []struct {
		// Input
		name   string
		chatDB ChatDBMock

		// Output
		messages []shared.Message
		err      error
	}{
		{
			// Input
			name: "Test1",
			chatDB: ChatDBMock{
				chatId: 0,
				exist:  true,
				messages: []shared.Message{
					shared.NewMessage(321, "hello"),
					shared.NewMessage(321, "how are u?"),
					shared.NewMessage(123, "im fine thanks"),
				},
				err: nil,
			},

			// Output
			err: error(nil),
			messages: []shared.Message{
				shared.NewMessage(321, "hello"),
				shared.NewMessage(321, "how are u?"),
				shared.NewMessage(123, "im fine thanks"),
			},
		},
		{
			// Input
			name: "Test2",
			chatDB: ChatDBMock{
				chatId:   0,
				exist:    false,
				messages: nil,
				err:      nil,
			},

			// Output
			messages: nil,
			err:      fmt.Errorf("chat with ID %d does not exist", chatID)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			c, err := GetChat(&tc.chatDB, 123)

			if err != tc.err {
				t.Errorf("Error while creating new chat: %v", err)
			} else {
				if c == nil {
					return
				}
			}

			messages, err := c.GetMessages()
			if err != nil {
				t.Errorf("Could not get messages: %v", err)
			}

			for i, msg := range messages {
				if msg.Text != tc.messages[i].Text {
					t.Errorf("Message text not equal: got %v, want %v", msg.Text, tc.messages[i].Text)
				}
				if msg.UserID != tc.messages[i].UserID {
					t.Errorf("Message UserID not equal: got %v, want %v", msg.UserID, tc.messages[i].UserID)
				}
			}
		})
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
