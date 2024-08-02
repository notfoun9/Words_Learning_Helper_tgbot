This is a bot that helps ypu to learn new English words!

To add a new word to your list send to to bot and then it will ask you for a definitoin. Send it and word+definition will be sawed in your list!
Command /rnd makes the bot send you a random word from your list and it's hidden definition. Tap of the definition text to reveal it and check if you guessed it right!
Command /all sends you all your saved words with hidden definitions.
Command /rmv makes the bot remove a word from your list. You will be asked for a word to remove after sending this command.

To lauch your own bot you need to ask a BotFather to create a new bot send you it's unique token.
After that execute command "go build | ./telegram-bot -tg-bot-token 'YOUR TOKEN'" from the main directory and it will work.

This bot uses external SQLite Go package to contain data of users and sends http request and responses to communicate with 'api.telegram.org' via standart Go "http" package. 

