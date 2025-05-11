package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"


func BlankMessage(update tgbotapi.Update) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "")
}

func NewStickerURL(update tgbotapi.Update, url string) tgbotapi.StickerConfig {
	file := tgbotapi.FileID(url)
	return tgbotapi.NewSticker(update.Message.Chat.ID, file)
}