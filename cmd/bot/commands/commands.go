package commands

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Reg(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	
	/*
	TODO: Solve the the full_name problem!
	*/

	msg := BlankMessage(update)
	
	chat_id := update.Message.Chat.ID
	reg_member_id := update.Message.From.ID
	
	user_name := update.Message.From.UserName
	full_name := update.Message.From.String()

	if user_name == full_name {
		full_name = ""
	}

	ok := CreateUser(chat_id, reg_member_id, full_name, user_name)

	var user_string string
	if full_name == "" {
		user_string = user_name
	} else {
		user_string = full_name
	}

	if ok {
		msg.Text = fmt.Sprintf("%s, te-ai inregistrat cu succes", user_string)
	} else {
		msg.Text = fmt.Sprintf("%s, dece te inregistrezi de 2 ori? 🤡", user_string)
	}

	_ , err := bot.Send(msg)

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

}