package db

type ChatDB interface {
	Store(msg string) error
	GetMessages() ([]string, error)
}

type ChatDBImpl struct {
	id1 uint32
	id2 uint32
}
