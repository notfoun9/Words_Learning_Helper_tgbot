package telegram_events

const msgHelp = `This bot helps to learn languages.
You can send a word and then you will be asked for a definition.
Use command /rnd to get a random word from the list and check if you remember it.
Use command /rmv to be asked for a word to be removed from the list.
Use command /all to get all the list printed`

const msgHello = "Hello there!\n\n" + msgHelp

const (
	msgGiveDefinition = "Give me a definition of this word 😳"
	msgUnknownCmd     = "Unknown command"
	msgNoSavedWords   = "You have no saved words left 🤡"
	msgSaved          = "Saved! ✅"
	msgAlreadyExists  = "You have already saved this word 🥴"
	msgListOfWords    = "Here are all your saved words: "
	msgNoSuchWord     = "You haven't saved this word 🤡"
	msgWordToDelete   = "Which word do you want to delete? 🥺"
	msgWordRemoved    = " is removed from the list 💀"
)
