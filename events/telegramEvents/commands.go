package telegram_events

import (
	"errors"
	"log"
	"os"
	"strings"
	"telegram-bot/clients/telegram_client"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd    = "/rnd"
	HelpCmd   = "/help"
	StartHelp = "/start"
	RemoveCmd = "/rmv"
	AllCmd    = "/all"

	ReadyToRemove = "removeWord"
)

func (proc *ProcessorTelegram) doCommand(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new cmd %s from %s", text, username)

	state, err := proc.wordStorage.GetState(username)
	if err != nil {
		return err
	}

	if state == ReadyToRemove {
		return proc.removeWord(chatID, username, text)
	} else if len(state) == 0 {
		switch text {
		case RndCmd:
			return proc.sendRandomWord(chatID, username)
		case HelpCmd:
			return proc.sendHelp(chatID)
		case StartHelp:
			return proc.sayHello(chatID)
		case RemoveCmd:
			return proc.removeCmd(chatID, username)
		case AllCmd:
			return proc.printAll(chatID, username)
		default:
			return proc.saveWord(chatID, text, username)
		}
	} else {
		return proc.giveDefinitionWord(chatID, username, text)
	}
}

func (proc *ProcessorTelegram) sendHelp(chatID int) error {
	return proc.tgClient.SendMessage(chatID, msgHelp)
}

func (proc *ProcessorTelegram) sayHello(chatID int) error {
	return proc.tgClient.SendMessage(chatID, msgHello)
}

func (proc *ProcessorTelegram) saveWord(chatID int, word string, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: save page", err) }()
	sendMessage := NewMessageSender(chatID, proc.tgClient)

	wordToAdd := &storage.Word{
		Word:     word,
		UserName: username,
	}

	doesExist, err := proc.wordStorage.DoesExist(username, word)
	if err != nil {
		return err
	}
	if doesExist {
		return sendMessage(msgAlreadyExists)
	}

	if err := proc.wordStorage.Save(wordToAdd); err != nil {
		return err
	}

	return sendMessage(msgGiveDefinition)
}

func (proc *ProcessorTelegram) sendRandomWord(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do comand: pick random", err) }()
	sendMessage := NewMessageSender(chatID, proc.tgClient)

	word, err := proc.wordStorage.PickRandom(username)

	if err != nil {
		if errors.Is(err, storage.ErrNoWordsSaved) {
			return sendMessage(msgNoSavedWords)
		} else {
			return err
		}
	}

	if err := sendMessage(word.Word); err != nil {
		return err
	}

	return proc.tgClient.SendSpoilerMessage(chatID, word.Definition)
}

func (proc *ProcessorTelegram) removeCmd(chatID int, username string) error {
	if err := proc.tgClient.SendMessage(chatID, msgWordToDelete); err != nil {
		return err
	}

	return proc.wordStorage.SetState(username, ReadyToRemove)
}

func (proc *ProcessorTelegram) removeWord(chatID int, username string, word string) error {
	err := proc.wordStorage.SetState(username, "")
	if err != nil {
		return err
	}

	err = proc.wordStorage.Remove(username, word)
	if err == os.ErrNotExist {
		return proc.tgClient.SendMessage(chatID, msgNoSuchWord)
	} else if err != nil {
		return err
	}

	return proc.tgClient.SendMessage(chatID, word+msgWordRemoved)
}

func (proc *ProcessorTelegram) giveDefinitionWord(chatID int, username string, def string) (err error) {
	err = proc.wordStorage.GiveDefinition(username, def)
	if err != nil {
		return err
	}
	return proc.tgClient.SendMessage(chatID, msgSaved)
}

func (proc *ProcessorTelegram) printAll(chatID int, username string) (err error) {
	defer func() { e.Wrap("Unable to print the list ", err) }()

	words, err := proc.wordStorage.AllWords(username)
	if err == storage.ErrNoWordsSaved {
		return proc.tgClient.SendMessage(chatID, msgNoSavedWords)
	} else if err != nil {
		return err
	}

	for i := 0; i < len(words); i++ {
		err := proc.tgClient.SendMessage(chatID, words[i].Word)
		if err != nil {
			return err
		}

		err = proc.tgClient.SendSpoilerMessage(chatID, words[i].Definition)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMessageSender(chatID int, tgClient *telegram_client.Client) func(string) error {
	return func(msg string) error {
		return tgClient.SendMessage(chatID, msg)
	}
}
