package main

import (
	"strings"

	"github.com/JojiiOfficial/gaw"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *Bot) isCommand(input, dest string) bool {
	dests := []string{"!" + dest, "/" + dest, "!" + dest + "@" + bot.Self.UserName, "/" + dest + "@" + bot.Self.UserName}

	if !strings.HasPrefix(input, "!") && strings.HasPrefix(input, "@") {
		return false
	}

	return gaw.IsInStringArray(input, dests)
}

func getNameFromUser(user *tgbotapi.User) string {
	if len(user.UserName) > 0 {
		return user.UserName
	}

	return user.FirstName + "" + user.LastName
}
