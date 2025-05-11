package commands

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI
}

func (u Update) getFullName() string {

	user := u.Update.Message.From

	full_name := user.FirstName
	if user.LastName != "" {
		full_name += " " + user.LastName
	}

	return full_name

}

func (u Update) pingText() string {

	ping_text := fmt.Sprintf("[@%s](tg://user?id=%d)", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, u.Update.Message.From.String()) , u.Update.Message.From.ID)

	// ping_text = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, ping_text)

	return ping_text
}

func (u Update) pingMessage(text string) tgbotapi.MessageConfig {

	msg := BlankMessage(u.Update)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	msg.Text = fmt.Sprintf("%s%s", u.pingText(), tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2 ,text))

	return msg
}

func (u Update) Reg() {

	var msg tgbotapi.MessageConfig

	chat_id := u.Update.Message.Chat.ID
	reg_member_id := u.Update.Message.From.ID

	user_name := u.Update.Message.From.UserName

	full_name := u.getFullName()

	//old
	// full_name := Update.Message.From.String()
	ok := CreateUser(chat_id, reg_member_id, full_name, user_name)

	// var user_string string
	// if full_name == "" {
	// 	user_string = user_name
	// } else {
	// 	user_string = full_name
	// }

	if ok {
		msg = u.pingMessage(", te-ai inregistrat cu succes")
	} else {
		msg= u.pingMessage(", dece te inregistrezi de 2 ori? 🤡")
	}

	_, err := u.Bot.Send(msg)

	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

}

func (u Update) Unreg() {
	
	chat_id := u.Update.Message.Chat.ID
	member_id := u.Update.Message.From.ID

	ok, err := DeleteUser(chat_id, member_id)
	
	if err != nil {
		log.Println("[ERROR]", err)
	}

	msg := BlankMessage(u.Update)

	if ok {
		msg = u.pingMessage(" a iesit cu pozor, dar statistica tine minte tot")  
	} else if err != nil {
		msg.Text = "utilizatorul nu a fost gasit"
	} else {
		msg.Text = fmt.Sprint("Internal Error ", err)
	}

	u.Bot.Send(msg)

}

func (u Update) EchoNickName() {

	member_id := u.Update.Message.From.ID

	msg := BlankMessage(u.Update)

	nick_name := GetNickName(member_id)

	if nick_name != "" {
		msg.Text = nick_name
	} else {
		msg.Text = "No nick_name found!"
	}

	u.Bot.Send(msg)

}

func (u Update) PingMe() {

	msg := BlankMessage(u.Update)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	ping_text := u.pingText()

	ping_text = fmt.Sprintf("I am pinging you, %s", ping_text)

	msg.Text = ping_text

	log.Println("[DEBUG] ", ping_text)

	u.Bot.Send(msg)

}