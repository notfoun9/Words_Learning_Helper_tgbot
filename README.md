This is a bot that helps you learn new English words!
Inspired by youtube.com/@nikolay_tuzov guide

To add a new word to your list send it to the bot and then it will ask you for a definitoin. Send it and the word+definition will be saved in your list!
The /rnd command makes the bot send you a random word from your list and its hidden definition. Click on the definition text to open it and check if you guessed correctly!
The /all command sends you all your saved words with their hidden definitions.
The /rmv command makes the bot remove a word from your list. You will be asked for a word to remove after sending this command.

To launch your own bot you need to ask a BotFather to create a new bot and send you its unique token.
After that run the command "go build | ./telegram-bot -tg-bot-token 'YOUR TOKEN'" from the main directory and the bot will start working.

This bot uses the external SQLite Go package to contain data of users and sends http requests and responses to communicate with 'api.telegram.org' via the standart Go "http" package. 

