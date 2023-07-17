package shared

import (
	"fmt"
	"time"
)

// IncomingMessage Messages coming from client to server
type IncomingMessage struct {
	ChatID   uint32 `json:"chat_id"`
	SenderID uint32 `json:"sender_id"`
	Text     string `json:"text"`
}

type Message struct {
	ID        uint32    `json:"id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	ChatID    uint32    `json:"chat_id"`
	SenderID  uint32    `json:"sender_id"`
}

func NewMessage(senderID uint32, chatID uint32, text string) (Message, error) {
	err := checkText(text)
	if err != nil {
		return Message{}, err
	}
	err = checkSenderID(senderID)
	if err != nil {
		return Message{}, err
	}

	return Message{
		SenderID:  senderID,
		ChatID:    chatID,
		Text:      text,
		Timestamp: time.Now(),
	}, nil
}

func NewIncomingMessage(senderID uint32, chatID uint32, text string) (IncomingMessage, error) {
	err := checkText(text)
	if err != nil {
		return IncomingMessage{}, err
	}
	err = checkSenderID(senderID)
	if err != nil {
		return IncomingMessage{}, err
	}

	return IncomingMessage{
		SenderID: senderID,
		ChatID:   chatID,
		Text:     text,
	}, nil
}

func checkText(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("empty message")
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
