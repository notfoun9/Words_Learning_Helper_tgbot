package storage

import (
	"errors"
)

var ErrNoWordsSaved = errors.New("no words saved")

type WordStorage interface {
	Save(word *Word) error
	PickRandom(username string) (*Word, error)
	Remove(username string, word string) error
	DoesExist(username string, word string) (bool, error)
	GiveDefinition(username string, definition string) error

	GetState(username string) (string, error)
	SetState(username string, state string) error

	AllWords(username string) ([]Word, error)
}

type Word struct {
	Word       string
	Definition string
	UserName   string
}
