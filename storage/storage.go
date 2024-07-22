package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"telegram-bot/lib/e"
)

var ErrNoPagesSaved = errors.New("no pages saved")

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	DoesExist(page *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
