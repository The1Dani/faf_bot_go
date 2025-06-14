package main

import (
	"log"
	"os"
	"strings"

	"github.com/The1Dani/faf_bot_go/cmd/bot/commands"
	// "github.com/The1Dani/faf_bot_go/cmd/bot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTelegramBot() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	upd := commands.Update{Bot: bot}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// // CHANGE AFTER USE
	// m := tgbotapi.NewMessage(messages.THE_SERVER, "Hello From rewritten Bot")
	// bot.Send(m)

	for update := range updates {

		upd.Update = update

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			reponse := strings.Split(update.CallbackQuery.Data, ":")

			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			switch reponse[0] {
			case commands.Carmic:
				upd.CallBackCarmic(reponse[1])
			}

		}

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "_")

		// Extract the command from the Message.

		switch update.Message.Command() {
		case "reg":
			upd.Reg()
		case "unreg":
			upd.Unreg()
		case "pidor":
			upd.Pidor()
		case "run":
			upd.Nice()
		case "stats":
			upd.Stats()
		case "pidorstats":
			upd.PidorStats()
		case "percentstats":
			upd.PercentStats()
		case "carmicdices":
			upd.Carmic()
		default:
			msg.Text = "I don't know that command"
			bot.Send(msg)
		}

		if err != nil {
			log.Panic(err)
		}

	}
}
