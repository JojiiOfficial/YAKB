package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	conf, err := initConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := NewDB(conf.DataFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if conf.AllowSelfKarma {
		fmt.Println("Self carma is allowed!")
	}

	// Register new TG bot
	tgbot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot := Bot{
		DatabaseHandler: db,
		BotAPI:          tgbot,
		config:          conf,
	}

	fmt.Printf("Bot %s registered\n", bot.Self.UserName)

	// start our bot logic
	err = bot.start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot started")

	bot.await()
}
