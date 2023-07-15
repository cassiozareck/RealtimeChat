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

	message1 := shared.IncomingMessage{ChatID: 123, SenderID: 12, Text: "Hi"}
	message2 := shared.IncomingMessage{ChatID: 123, SenderID: 15, Text: "Goodbye!"}

	var testCases = []struct {
		// Input
		name     string
		messages []shared.IncomingMessage
		chatDB   db.ChatDB

		// Output
		outMessages []shared.Message
		outPeople   []uint32
		err         error
	}{
		{
			name:     "Test 1",
			messages: []shared.IncomingMessage{message1, message2},
			chatDB:   &ChatDBMock{},

			outMessages: []shared.Message{{
				ID:       0,
				SenderID: message1.SenderID,
				ChatID:   message1.ChatID,
				Text:     message1.Text,
			}, {
				ID:       0,
				SenderID: message2.SenderID,
				ChatID:   message2.ChatID,
				Text:     message2.Text,
			}},
			outPeople: []uint32{12, 15},
			err:       nil,
		},
		{
			name:     "Test 2",
			messages: []shared.IncomingMessage{message2},
			chatDB: &ChatDBMock{
				Messages: []shared.Message{{
					ID:       0,
					SenderID: message1.SenderID,
					ChatID:   message1.ChatID,
					Text:     message1.Text,
				},
				}},
			outMessages: []shared.Message{{
				ID:       0,
				SenderID: message1.SenderID,
				ChatID:   message1.ChatID,
				Text:     message1.Text,
			}, {
				ID:       0,
				SenderID: message2.SenderID,
				ChatID:   message2.ChatID,
				Text:     message2.Text,
			}},
			outPeople: []uint32{message1.SenderID, message2.SenderID},
			err:       nil,
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

			chatMessages, err := c.GetMessages()
			if err != nil {
				t.Fatal("Error while getting messages: ", err)
			}

			// We don't check time because it is not set
			for i, msg := range tc.outMessages {
				if msg.ID != chatMessages[i].ID ||
					msg.ChatID != chatMessages[i].ChatID ||
					msg.SenderID != chatMessages[i].SenderID ||
					msg.Text != chatMessages[i].Text {
					t.Fatalf("Expected %v, got %v", msg, chatMessages[i])
				}
			}

			chatPeople, err := c.GetPeople()
			if err != nil {
				t.Fatal("Error while getting people: ", err)
			}
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
