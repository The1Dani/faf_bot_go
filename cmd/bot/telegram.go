package main

import (
	"log"
	"os"

	"github.com/The1Dani/faf_bot_go/cmd/bot/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartTelegramBot() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
    if err != nil {
        log.Panic(err)
    }

    upd := commands.Update{Bot: bot}

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {

        upd.Update = update

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
        case "help":
            msg.Text = "I understand /sayhi and /status."
        case "sayhi":
            msg.Text = "Hi :)"
        case "status":
            msg.Text = "I'm ok."
        case "reg":
            upd.Reg()
        case "nick":
            upd.EchoNickName()
        case "unreg":
            upd.Unreg()
        case "pingme":
            upd.PingMe()
        case "sticker":
            upd.SendSticker()

        default:
            msg.Text = "I don't know that command"
        }

		_, err := bot.Send(msg)

        if err != nil {
            log.Panic(err)
        }

		log.Println(update.Message.From.FirstName)
		
    }
}
