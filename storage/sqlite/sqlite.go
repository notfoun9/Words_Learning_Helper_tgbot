package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"telegram-bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPermission = 0774
const databaseName = "sqliteDB"
const state = "statestatestate"

type SqlStorage struct {
	dataBase *sql.DB
}

func NewSQLStorage(folderPath string) (*SqlStorage, error) {
	if err := os.MkdirAll(folderPath, defaultPermission); err != nil {
		return nil, err
	}

	dataBase, err := sql.Open("sqlite3", filepath.Join(folderPath, databaseName))
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := dataBase.Ping(); err != nil {
		return nil, fmt.Errorf("can't conect to database: %w", err)
	}

	return &SqlStorage{
		dataBase: dataBase,
	}, nil
}

func (sq *SqlStorage) Init() error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS words (username TEXT, word TEXT, definition TEXT)`

	_, err := sq.dataBase.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("can't create a table: %w", err)
	}
	return nil
}

func (sq SqlStorage) Save(word *storage.Word) error {
	saveWorsQuery := `INSERT INTO words (username, word, definition) VALUES(?, ?, ?)`
	_, err := sq.dataBase.Exec(saveWorsQuery, word.UserName, word.Word, "")
	if err != nil {
		return fmt.Errorf("cant save the word whaaat: %w", err)
	}

	sq.SetState(word.UserName, word.Word)
	return nil
}

func (sq *SqlStorage) PickRandom(username string) (*storage.Word, error) {
	size, err := sq.size(username)
	if err != nil {
		return nil, err
	}
	if size == 1 {
		return nil, storage.ErrNoWordsSaved
	}

	queryGetWord := `SELECT word FROM words WHERE username = ? AND word != ? ORDER BY RANDOM() LIMIT 1`
	queryGetDefinition := `SELECT definition FROM words WHERE username = ? AND word = ? ORDER BY RANDOM() LIMIT 1`

	var word string
	var definition string

	err = sq.dataBase.QueryRow(queryGetWord, username, state).Scan(&word)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoWordsSaved
	} else if err != nil {
		return nil, fmt.Errorf("cant save the word: %w", err)
	}

	err = sq.dataBase.QueryRow(queryGetDefinition, username, word).Scan(&definition)
	if err != nil {
		return nil, fmt.Errorf("cant save the word: %w", err)
	}
	return &storage.Word{
		Word:       word,
		Definition: definition,
		UserName:   username,
	}, nil
}

func (sq *SqlStorage) Remove(username string, word string) error {
	doesExist, err := sq.DoesExist(username, word)
	if err != nil {
		return err
	}
	if !doesExist {
		return os.ErrNotExist
	}

	queryDeleteWord := `DELETE FROM words WHERE username = ? AND word = ?`
	if _, err := sq.dataBase.Exec(queryDeleteWord, username, word); err != nil {
		return fmt.Errorf("cant remove page: %w", err)
	}
	return nil
}

func (sq *SqlStorage) DoesExist(username string, word string) (bool, error) {
	queryCountEntries := `SELECT COUNT(*) FROM words WHERE username = ? AND word = ?`

	var count int

	err := sq.dataBase.QueryRow(queryCountEntries, username, word).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("cant check if the word is saved: %w", err)
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (sq *SqlStorage) GiveDefinition(username string, definition string) error {
	querySetDefinition := `UPDATE words SET definition = ? WHERE username = ? AND word = ?;`

	word, err := sq.GetState(username)
	if err != nil {
		return err
	}

	_, err = sq.dataBase.Exec(querySetDefinition, definition, username, word)
	if err != nil {
		return fmt.Errorf("cant give the definition: %w", err)
	}

	sq.SetState(username, "")
	return nil
}

func (sq *SqlStorage) GetState(username string) (string, error) {
	doesExist, err := sq.DoesExist(username, state)
	if err != nil {
		log.Println("impossible to get state: ", err)
		return "nil", err
	}

	if !doesExist {
		sq.Save(&storage.Word{
			Word:       state,
			Definition: "",
			UserName:   username,
		})
		return "", sq.SetState(username, "")
	}
	queryGetState := `SELECT definition FROM words WHERE username = ? AND word = ? ORDER BY RANDOM() LIMIT 1`

	var stateMsg string

	err = sq.dataBase.QueryRow(queryGetState, username, state).Scan(&stateMsg)
	if err != nil {
		return "", fmt.Errorf("cant get the state: %w", err)
	}

	return stateMsg, nil
}

func (sq *SqlStorage) SetState(username string, newState string) error {
	query := `UPDATE words SET definition = ? WHERE username = ? AND word = ?;`

	_, err := sq.dataBase.Exec(query, newState, username, state)
	if err != nil {
		return fmt.Errorf("cant set the new state: %w", err)
	}

	return nil
}

func (sq *SqlStorage) AllWords(username string) (words []storage.Word, err error) {
	size, err := sq.size(username)
	if err != nil {
		return nil, err
	}
	if size == 1 {
		return nil, storage.ErrNoWordsSaved
	}

	query := `SELECT word, definition FROM words WHERE username = ?`
	rows, err := sq.dataBase.Query(query, username)
	if err != nil {
		return nil, err
	}

	var word string
	var definition string

	for rows.Next() {
		err = rows.Scan(&word, &definition)
		if err != nil {
			return nil, err
		}
		if word != state {
			words = append(words, storage.Word{
				Word:       word,
				Definition: definition,
				UserName:   username,
			})
		}
	}

	return words, nil
}

func (sq *SqlStorage) size(username string) (size int, err error) {
	query := `SELECT COUNT(username) FROM words WHERE username = ?`

	err = sq.dataBase.QueryRow(query, username).Scan(&size)
	return
}
