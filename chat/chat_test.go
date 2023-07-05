package chat

import (
	"github.com/cassiozareck/realchat/db"
	"github.com/cassiozareck/realchat/shared"
	"testing"
)

type ChatDBMock struct {
	chatId           uint32
	exist            bool
	Messages         []shared.Message
	err              error
	storeErr         error
	getMessagesError error
}

func (c *ChatDBMock) CreateChat() (uint32, error) {
	return c.chatId, c.err
}

func (c *ChatDBMock) ChatExists(chatID uint32) (bool, error) {
	return c.exist, c.err
}

func (c *ChatDBMock) Store(msg shared.Message) error {
	c.Messages = append(c.Messages, msg)
	return c.storeErr
}

func (c *ChatDBMock) GetMessages(chatID uint32) ([]shared.Message, error) {
	return c.Messages, c.getMessagesError
}

// TestChat will send a message and test if the
// last message is giving the appropriate answer
func TestChat(t *testing.T) {

	message1, _ := shared.NewMessage(123, "Hi")
	message2, _ := shared.NewMessage(321, "Goodbye!")

	var testCases = []struct {
		// Input
		name     string
		messages []shared.Message
		chatDB   db.ChatDB

		// Output
		outMessages []shared.Message
		outPeople   []uint32
		err         error
	}{
		{
			name:     "Test 1",
			messages: []shared.Message{message1, message2},
			chatDB:   &ChatDBMock{},

			outMessages: []shared.Message{message1, message2},
			outPeople:   []uint32{123, 321},
			err:         nil,
		},
		{
			name:     "Test 2",
			messages: []shared.Message{message2},
			chatDB: &ChatDBMock{
				Messages: []shared.Message{message1},
			},

			outMessages: []shared.Message{message1, message2},
			outPeople:   []uint32{123, 321},
			err:         nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewChat(tc.chatDB)
			if err != nil {
				t.Fatal("Error while creating new chat: ", err)
			}

			for _, msg := range tc.messages {
				err = c.SendMessage(msg)
				if err != nil {
					t.Fatal("Error while sending message: ", err)
				}
			}

			err = c.UpdateMessages()
			if err != nil {
				t.Fatal("Error while updating messages: ", err)
			}

			chatMessages := c.GetMessages()

			for i, msg := range tc.outMessages {
				if msg != chatMessages[i] {
					t.Fatalf("Expected %v, got %v", msg, chatMessages[i])
				}
			}

			err = c.UpdatePeople()
			if err != nil {
				t.Fatal("Error while updating people: ", err)
			}

			chatPeople := c.GetPeople()
			for i, person := range tc.outPeople {
				if person != chatPeople[i] {
					t.Fatalf("Expected %v, got %v", person, chatPeople[i])
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

func TestChat_WithError(t *testing.T) {
	var testCases = []struct {
		// Input
		name   string
		chatDB db.ChatDB

		// Output
		err error
	}{
		{
			name: "Test 1",
			chatDB: &ChatDBMock{
				storeErr: error(nil),
			},
			err: error(nil),
		},
		{
			name: "Test 2",
			chatDB: &ChatDBMock{
				getMessagesError: error(nil),
			},
			err: error(nil),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewChat(tc.chatDB)
			if err != tc.err {
				t.Fatalf("Expected %v, got %v", tc.err, err)
			}
		})
	}
}
