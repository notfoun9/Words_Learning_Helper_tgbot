package storage

import (
	"errors"
)

var ErrNoWordsSaved = errors.New("no words saved")

type WordStorage interface {
	SaveWord(word *Word) error
	PickRandomWord(username string) (*Word, error)
	RemoveWord(username string, word string) error
	DoesExistWord(username string, word string) (bool, error)
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
