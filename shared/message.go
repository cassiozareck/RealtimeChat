package shared

import (
	"fmt"
	"time"
)

type Message struct {
	id        uint32
	chatID    uint32
	senderID  uint32
	text      string
	timestamp time.Time
}

func NewMessage(senderID uint32, text string) (Message, error) {
	err := checkText(text)
	if err != nil {
		return Message{}, err
	}
	err = checkSenderID(senderID)
	if err != nil {
		return Message{}, err
	}

	return Message{
		senderID:  senderID,
		text:      text,
		timestamp: time.Now(),
	}, nil
}

func NewMessageFromDB(id uint32, chatID uint32, senderID uint32, text string, timestamp time.Time) Message {
	return Message{
		id:        id,
		chatID:    chatID,
		senderID:  senderID,
		text:      text,
		timestamp: timestamp,
	}
}

func (m *Message) ID() uint32 {
	return m.id
}

func (m *Message) ChatID() uint32 {
	return m.chatID
}

func (m *Message) SenderID() uint32 {
	return m.senderID
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

func (m *Message) SetText(text string) error {
	err := checkText(text)
	if err != nil {
		return err
	}
	m.text = text
	return nil
}

func checkText(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("invalid message: %s", text)
	}
	if len(text) > MAX_MESSAGE_SIZE {
		return fmt.Errorf("message too long: %s", text)
	}
	return nil
}

func checkSenderID(SenderID uint32) error {
	if SenderID <= 0 {
		return fmt.Errorf("invalid sender id: %d", SenderID)
	}
	return nil
}
