package shared

import "time"

type Message struct {
	ID     uint32
	ChatID uint32
	UserID uint32
	Text   string
	Hour   time.Time
}

func (m Message) LessThan(other Message) bool {
	return m.Hour.Before(other.Hour)
}

func NewMessage(userID uint32, text string) Message {
	return Message{
		UserID: userID,
		Text:   text,
		Hour:   time.Now(),
	}
}
