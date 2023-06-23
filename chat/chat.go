package chat

import "github.com/cassiozareck/realchat/db"

type Chat struct {
	chatDB db.ChatDB
}

func GetChat(chatDB db.ChatDB) *Chat {
	c := Chat{chatDB}
	return &c
}

func (t *Chat) SendMessage(msg string) error {
	err := t.chatDB.Store(msg)
	if err != nil {
		return err
	}
	return nil
}

func (t *Chat) GetMessages() ([]string, error) {
	msgs, err := t.chatDB.GetMessages()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (t *Chat) LastMessage() (string, error) {
	msgs, err := t.chatDB.GetMessages()
	if err != nil {
		return "", err
	}
	return msgs[len(msgs)-1], nil
}
