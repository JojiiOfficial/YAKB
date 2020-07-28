package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/JojiiOfficial/gaw"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	clitable "gopkg.in/benweidig/cli-table.v2"
)

// Bot is this bot lol
type Bot struct {
	*tgbotapi.BotAPI
	*DatabaseHandler

	config            *Config
	exit              chan struct{}
	lastChatMessageID map[int64]int // map[chatID]MessageID
	mx                sync.Mutex
}

func (bot *Bot) start() error {
	bot.lastChatMessageID = make(map[int64]int)
	bot.exit = make(chan struct{}, 1)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			go bot.handleMessage(update)
		}
	}()

	return nil
}

func (bot *Bot) handleMessage(update tgbotapi.Update) {
	messageText := update.Message.Text

	isReply := update.Message.ReplyToMessage != nil

	if messageText == bot.config.KarmaTopCommand {
		var id int
		if isReply {
			id = update.Message.ReplyToMessage.From.ID
		}

		msgTxt := bot.getKarmaTop(id)
		if len(msgTxt) > 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTxt)
			bot.Send(msg)
		}

		return
	}

	// Ignore non replies
	if !isReply {
		return
	}

	// Ignore other messages
	if bot.ignoreMessage(messageText) {
		return
	}

	// Prevent selfecarma if disabled
	if !bot.config.AllowSelfKarma &&
		update.Message.ReplyToMessage.From.ID == update.Message.From.ID {
		return
	}

	var karmaDelta int
	if bot.isKarmaAdd(messageText) {
		karmaDelta = 1
	}

	if bot.isKarmaRemove(messageText) {
		karmaDelta = -1
	}

	replyTo := update.Message.ReplyToMessage
	succ, err := bot.addKarma(replyTo.MessageID, update.Message.From.ID, replyTo.From, karmaDelta)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !succ {
		return
	}

	bot.mx.Lock()
	bot.runNotificationHook(update, karmaDelta)
	bot.mx.Unlock()

	if bot.isKarmaRemove(messageText) && replyTo.From.ID == bot.Self.ID {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "kek")
		msg.ReplyToMessageID = update.Message.MessageID

		r, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}

		bot.handleMessage(tgbotapi.Update{
			Message: &r,
		})
	}
}

func (bot *Bot) runNotificationHook(update tgbotapi.Update, kDelta int) {
	cid := update.Message.Chat.ID

	fmt.Println(bot.lastChatMessageID)

	if msg, has := bot.lastChatMessageID[cid]; has {
		dm := tgbotapi.NewDeleteMessage(cid, msg)
		_, err := bot.DeleteMessage(dm)
		if err != nil {
			fmt.Println(err)
		}
	}
	user := update.Message.ReplyToMessage.From

	txt := "Karma"
	if kDelta > 0 {
		txt += " incremented by " + strconv.Itoa(int(math.Abs(float64(kDelta))))
	} else {
		txt += " decremented by " + strconv.Itoa(int(math.Abs(float64(kDelta))))
	}
	txt += " for @" + user.UserName

	msg := tgbotapi.NewMessage(cid, txt)
	r, err := bot.Send(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.lastChatMessageID[cid] = r.MessageID
}

func (bot *Bot) getKarmaTop(userid int) string {
	karmas, err := bot.getTopKarma(userid)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	table := clitable.New()
	table.ColSeparator = " "
	table.Padding = 1

	for i := range karmas {
		table.AddRow(karmas[i].UserName, karmas[i].Amount)
	}

	return table.String()
}

func (bot *Bot) isKarmaAdd(message string) bool {
	return gaw.IsInStringArray(strings.ToLower(message), bot.config.AddKarma)
}

func (bot *Bot) isKarmaRemove(message string) bool {
	return gaw.IsInStringArray(strings.ToLower(message), bot.config.RemoveKarma)
}

func (bot *Bot) ignoreMessage(message string) bool {
	return !bot.isKarmaAdd(message) && !bot.isKarmaRemove(message)
}

func (bot *Bot) await() {
	<-bot.exit
}
