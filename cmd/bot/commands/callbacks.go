package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	Carmic = "carmic"
)

func (u Update) CallBackCarmic(val string) {
	call_message := u.Update.CallbackQuery.Message
	
	msg := tgbotapi.NewEditMessageText(call_message.Chat.ID, call_message.MessageID, "") 
	
	switch val {
		case "true":
			msg.Text = "Cuburile karmei sunt activate"
			SetCarmic(call_message.Chat.ID, true)
		case "false":
			msg.Text = "Cuburile karmei sunt dezactivate"
			SetCarmic(call_message.Chat.ID, false)
	}
	
	u.Bot.Send(msg)
} 
