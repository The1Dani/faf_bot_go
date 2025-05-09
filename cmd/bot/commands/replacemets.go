package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"


func BlankMessage(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "")
}