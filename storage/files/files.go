package files

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
	"time"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(base string) WordStorage {
	return WordStorage{basePath: base}
}

type WordStorage struct {
	basePath string
}

func NewWordsStorage(base string) WordStorage {
	return WordStorage{basePath: base}
}

func (ws WordStorage) SaveWord(word *storage.Word) (err error) {
	defer func() { err = e.Wrap("can't save the word", err) }()

	folderPath := filepath.Join(ws.basePath, word.UserName)
	if err := os.MkdirAll(folderPath, defaultPerm); err != nil {
		return err
	}

	fileName := word.Word + ".txt"

	filePath := filepath.Join(folderPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err != nil {
		return err
	}

	path := filepath.Join(ws.basePath, word.UserName, "StateFile.txt")

	if err := os.Truncate(path, 0); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(word.Word), defaultPerm); err != nil {
		return err
	}

	return nil
}

func (ws WordStorage) createStateFile(username string) error {
	folderPath := filepath.Join(ws.basePath, username)

	if err := os.MkdirAll(folderPath, defaultPerm); err != nil {
		return err
	}

	path := filepath.Join(folderPath, "StateFile.txt")

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		return err
	} else {
		return err
	}
}

func (ws WordStorage) PickRandomWord(username string) (word *storage.Word, err error) {
	defer func() { err = e.Wrap("cant pick a page", err) }()

	fPath := filepath.Join(ws.basePath, username)
	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 1 {
		return nil, storage.ErrNoPagesSaved
	}

	var file fs.DirEntry
	for {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(files))
		file = files[n]

		if file.Name() != "StateFile.txt" {
			break
		}
	}

	text, err := os.ReadFile(path.Join(fPath, file.Name()))
	if err != nil {
		return nil, err
	}

	w, _ := strings.CutSuffix(file.Name(), ".txt")

	return &storage.Word{
		Word:       w,
		Definition: string(text),
		UserName:   username,
	}, nil
}

func (ws WordStorage) RemoveWord(username string, word string) error {
	fileName := word + ".txt"

	path := filepath.Join(ws.basePath, username, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap("cant remove file", err)
	}

	path = filepath.Join(ws.basePath, username, "StateFile.txt")

	if err := os.Truncate(path, 0); err != nil {
		return err
	}

	return nil
}

func (ws WordStorage) DoesExistWord(username string, w string) (b bool, err error) {
	fileName := w + ".txt"

	path := filepath.Join(ws.basePath, username, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap("cant check if file exists", err)
	}

	return true, nil
}

func (ws WordStorage) GiveDefinition(username string, definition string) error {
	path := filepath.Join(ws.basePath, username, "StateFile.txt")
	word, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	wordPath := filepath.Join(ws.basePath, username, string(word)) + ".txt"

	if err := os.WriteFile(wordPath, []byte(definition), defaultPerm); err != nil {
		return err
	}

	if err := os.Truncate(path, 0); err != nil {
		return err
	}
	return nil
}

func (ws WordStorage) GetState(username string) (string, error) {
	if err := ws.createStateFile(username); err != nil {
		return "", err
	}
	path := filepath.Join(ws.basePath, username, "StateFile.txt")

	text, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func (ws WordStorage) SetState(username string, state string) error {
	path := filepath.Join(ws.basePath, username, "StateFile.txt")

	if err := os.Truncate(path, 0); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(state), defaultPerm)
}

func (ws WordStorage) AllWords(username string) ([]storage.Word, error) {
	fPath := filepath.Join(ws.basePath, username)
	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 1 {
		return nil, storage.ErrNoPagesSaved
	}

	var words []storage.Word
	for i := 0; i < len(files); i++ {
		if files[i].Name() == "StateFile.txt" {
			continue
		}

		def, err := os.ReadFile(path.Join(fPath, files[i].Name()))
		if err != nil {
			return nil, err
		}

		words = append(words, storage.Word{
			UserName:   username,
			Word:       strings.TrimSuffix(files[i].Name(), ".txt"),
			Definition: string(def),
		})
	}
	return words, nil
}
