package chat

import "testing"

type ChatDBMock struct {
	messages []string
}

func (c *ChatDBMock) Store(msg string) error {
	c.messages = append(c.messages, msg)
	return nil
}

func (c *ChatDBMock) GetMessages() ([]string, error) {
	return c.messages, nil
}

// TestSendMSG will send a message and test if the
// last message is giving the appropriate answer
func TestSendMSG(t *testing.T) {

	chatDB := ChatDBMock{}

	c := GetChat(&chatDB)

	MSG := "HELLO"
	err := c.SendMessage(MSG)
	if err != nil {
		t.Fatal("Error while sending message: ", err)
	}

	lastMessage, err := c.LastMessage()
	if err != nil {
		t.Fatal("Error while retrieving last message: ", err)
	}

	if lastMessage != MSG {
		t.Fatal("Last Message not equals ", MSG)
	}

}
