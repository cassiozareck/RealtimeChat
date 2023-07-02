package chat

import (
	"fmt"
	"github.com/cassiozareck/realchat/db"
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
	text := "hello"
	userID := uint32(2)

	var testCases = []struct {
		// Input
		name    string
		message string
		chatDB  db.ChatDB
		userID  uint32
		// Output
		outMessage shared.Message
		err        error
	}{
		{
			name:    "Test 1",
			message: text,
			chatDB: &ChatDBMock{
				chatId:   0,
				exist:    true,
				messages: []shared.Message{shared.NewMessage(userID, text)},
				err:      nil,
			},
			userID:     userID,
			outMessage: shared.NewMessage(userID, text),
			err:        nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewChat(tc.chatDB)
			if err != nil {
				t.Fatal("Error while creating new chat: ", err)
			}

			err = c.SendMessage(tc.userID, tc.message)
			if err != nil {
				t.Fatal("Error while sending message: ", err)
			}

			lastMessage, err := c.LastMessage()
			if err != nil {
				t.Fatal("Error while retrieving last message: ", err)
			}

			if lastMessage.Text != tc.outMessage.Text {
				t.Fatal("Last Message not equals ", tc.outMessage.Text)
			}

		})
	}
}

// TestChat_GetMessages will get messages using Chat and compare to see if
// messages are matching the same from database
func TestChat_GetMessages(t *testing.T) {
	chatID := uint32(2)
	testCases := []struct {
		// Input
		name   string
		chatDB ChatDBMock
		chatID uint32
		// Output
		messages []shared.Message
		err      error
	}{
		{
			// Input
			name:   "Test1",
			chatID: chatID,
			chatDB: ChatDBMock{
				exist: true,
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
			name:   "Test2",
			chatID: chatID,
			chatDB: ChatDBMock{
				exist:    false,
				messages: nil,
				err:      nil,
			},

			// Output
			messages: nil,
			err:      fmt.Errorf("chat with id %v does not exist", chatID)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			c, err := GetChat(&tc.chatDB, tc.chatID)

			if err != nil {
				if err.Error() != tc.err.Error() {
					t.Errorf("Error while creating new chat: %v", err)
				} else {
					if c == nil {
						return
					}
				}
			}

			messages := c.GetMessages()
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

func TestChat_GetPeople(t *testing.T) {
	personID1 := uint32(1)
	personID2 := uint32(2)

	testCases := []struct {
		// Input
		name   string
		chatDB ChatDBMock
		chatID uint32
		// Output
		people []shared.Person
		err    error
	}{
		{
			// Input
			name:   "Test1",
			chatID: 0,
			chatDB: ChatDBMock{
				exist: true,
				messages: []shared.Message{
					shared.NewMessage(personID1, "hello"),
					shared.NewMessage(personID1, "how are u?"),
					shared.NewMessage(personID2, "im fine thanks"),
				},
				err: nil,
			},

			// Output
			err: error(nil),
			people: []shared.Person{
				{
					ID:   personID1,
					Name: "",
				},
				{
					ID:   personID2,
					Name: "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			c, err := GetChat(&tc.chatDB, tc.chatID)

			if err != nil {
				t.Errorf("Error while creating new chat: %v", err)
			}

			people := c.GetPeople()
			if err != nil {
				t.Errorf("Could not get people: %v", err)
			}

			for i, person := range people {
				if person.ID != tc.people[i].ID {
					t.Errorf("Person ID not equal: got %v, want %v", person.ID, tc.people[i].ID)
				}
			}
		})
	}
}
