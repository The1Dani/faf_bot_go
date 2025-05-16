package commands

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/The1Dani/faf_bot_go/cmd/bot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmcvetta/randutil"
)

type Update struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI
}

type user struct {
	full_name string
	nick_name string
	member_id int64
	coefficient int32
	pidor_coefficient int32
}

type counts struct {
	PidorCount	int
	NiceCount	int
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

func (u user) pingText() string {

	user_string := u.nick_name;

	if user_string == "" {
		user_string = u.full_name
	}

	ping_text := fmt.Sprintf("[@%s](tg://user?id=%d)", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, user_string) , u.member_id)

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

func (u Update) Pidor() {

	chat_id := u.Update.Message.Chat.ID

	msg := BlankMessage(u.Update)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	congrats := BlankMessage(u.Update)
	sticker := NewStickerURL(u.Update, messages.BILLY_TEAR_OFF_VEST)


	ok, curr_user, curr_opp_user := TimeNotExpired(chat_id, pidor) // TEST | Add funct for getting the opposite modes user

	if ok {
		msg.Text = fmt.Sprintf("Pidorul zilei este deja selectet, este %s \\(%s\\)", curr_user.full_name, curr_user.pingText())	

	} else {
		
		var pidor_user user 
		var ok bool

		if CarmicDicesEnabled(chat_id) { // TEST
			ok, pidor_user = getRandomUserCarmic(chat_id, curr_opp_user, pidor) // TEST
		} else {
			ok, pidor_user = getRandomUser(chat_id, curr_opp_user) // TEST
		}

		if !ok {
			msg.Text = "Imposibil de selectat, lista de candidati e goala"
			u.Bot.Send(msg)
			return
		}

		pidor_count := UpdateStats(chat_id, pidor_user.member_id, pidor_stats) //TEST

		bdy := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("Pidorul zilei - %s ", pidor_user.full_name))

		msg.Text = fmt.Sprintf("%s %s", bdy, pidor_user.pingText())

		for _, txt := range messages.PIDOR_MESSAGES {
			pMsg := BlankMessage(u.Update)
			pMsg.Text = txt
			u.Bot.Send(pMsg)
			time.Sleep(1 * time.Second)
		}
		
		UpdateCurrent(chat_id, pidor_user.member_id, pidor)

		switch pidor_count {
			case 1:
				congrats.Text = messages.PIDOR_1_TIME
			case 10:
				congrats.Text = messages.PIDOR_10_TIME
			case 50:
				congrats.Text = messages.PIDOR_50_TIME
			case 100:
				congrats.Text = messages.PIDOR_100_TIME
		}
	}

	u.Bot.Send(msg)

	if congrats.Text != "" {
		u.Bot.Send(congrats)
		u.Bot.Send(sticker)
	}

}

func (u Update) Nice() {

	chat_id := u.Update.Message.Chat.ID

	msg := BlankMessage(u.Update)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	congrats := BlankMessage(u.Update)
	sticker := NewStickerURL(u.Update, messages.DRINK_CHAMPAGNE)

	ok, curr_user, curr_opp_user := TimeNotExpired(chat_id, nice) // TEST

	if ok {
		msg.Text = fmt.Sprintf("Krasavciku e deja selectat, e %s \\(%s\\)", curr_user.full_name, curr_user.pingText())	

	} else {
		
		var nice_user user 
		var ok bool

		if CarmicDicesEnabled(chat_id) { // TEST
			ok, nice_user = getRandomUserCarmic(chat_id, curr_opp_user, nice) // TEST
		} else {
			ok, nice_user = getRandomUser(chat_id, curr_opp_user) // TEST
		}

		if !ok {
			msg.Text = "Imposibil de selectat, lista de candidati e goala"
			u.Bot.Send(msg)
			return
		}

		nice_count := UpdateStats(chat_id, nice_user.member_id, nice_stats) //TEST

		bdy := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, fmt.Sprintf("Pidorul zilei - %s ", nice_user.full_name))

		msg.Text = fmt.Sprintf("%s %s", bdy, nice_user.pingText())

		for _, txt := range messages.NICE_MESSAGES {
			pMsg := BlankMessage(u.Update)
			pMsg.Text = txt
			u.Bot.Send(pMsg)
			time.Sleep(1 * time.Second)
		}
		
		UpdateCurrent(chat_id, nice_user.member_id, nice)

		switch nice_count {
			case 1:
				congrats.Text = messages.NICE_1_TIME
			case 10:
				congrats.Text = messages.NICE_10_TIME
			case 50:
				congrats.Text = messages.NICE_50_TIME
			case 100:
				congrats.Text = messages.NICE_100_TIME
		}
	}

	u.Bot.Send(msg)

	if congrats.Text != "" {
		u.Bot.Send(congrats)
		u.Bot.Send(sticker)
	}

}


func (u Update) EchoNickName() {

	member_id := u.Update.Message.From.ID
	chat_id := u.Update.Message.Chat.ID
	
	msg := BlankMessage(u.Update)

	user, _ := GetUser(member_id, chat_id)
	nick_name := user.nick_name
	
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

func (u Update) SendSticker() {

	chat_id := u.Update.Message.Chat.ID

	stF := tgbotapi.FileURL(messages.BILLY_TEAR_OFF_VEST)

	st := tgbotapi.NewSticker(chat_id, stF)

	u.Bot.Send(st)

}

func getRandomUserCarmic(chat_id int64, immune user, mode current_) (bool, user) {
	
	var members []user
	
	members, err := GetAllMembers(members, chat_id)
	
	if err != nil {
		return false, user{}
	}
	
	if len(members) < 1 {
		return false, user{}
	}
	
	members = slices.DeleteFunc(members, func(u user) bool {
		if u.member_id == immune.member_id {
			return true
		} else {
			return false
		}
	})
	
	if len(members) < 1 {
		return false, user{}
	}
	
	var choices []randutil.Choice
	
	for _, m := range members {
		c := randutil.Choice{}
		switch mode {
			case pidor:
				c.Weight = int(m.pidor_coefficient)
				c.Item = m
			
			case nice:
				c.Weight = int(m.coefficient)
				c.Item = m
			
			default:
			return false, user{}
		}
		choices = append(choices, c)
	}
	
	r, err := randutil.WeightedChoice(choices)
	
	if err != nil {
		log.Println("[ERROR]", err)
		return false, user{}
	}
	
	return true, r.Item.(user)
}

func getRandomUser(chat_id int64, immune user) (bool, user) {
	
	var members []user
	
	members, err := GetAllMembers(members, chat_id)
	
	if err != nil {
		return false, user{}
	}
	
	if len(members) < 1 {
		return false, user{}
	}
	
	members = slices.DeleteFunc(members, func(u user) bool {
		if u.member_id == immune.member_id {
			return true
		} else {
			return false
		}
	})
	
	if len(members) < 1 {
		return false, user{}
	}
	
	rand_member := members[rand.Intn(len(members))]
	
	return true, rand_member
}

func (u Update) Stats() {
	
	var text string
	
	msg := BlankMessage(u.Update)
	chat_id := u.Update.Message.Chat.ID
	results, members, err := GetStats(chat_id)
	text_list := []string{"Rezultatele jocului krasavciku zilei:"}

	
	if err == sql.ErrNoRows || err != nil {
		msg.Text = "Nimei nu e inregistrat, statistica e goala"
		u.Bot.Send(msg)
		u.Bot.Send(tgbotapi.NewMessage(chat_id, fmt.Sprint(err)))
		return
	}
	
	sort.Slice(members, func(i, j int) bool {
		return results[members[i].member_id].NiceCount < results[members[j].member_id].NiceCount
	})
	
	for _, memb := range members {
		txt := fmt.Sprintf("%v: %v", memb.full_name, results[memb.member_id].NiceCount)
		text_list = append(text_list, txt)
	}
	
	// log.Printf("[DEBUG] %#v\n%#v\n%#v\n", results, members, text_list)
	text = strings.Join(text_list, "\n")
	
	msg.Text = text
	u.Bot.Send(msg)
}

func (u Update) PidorStats() {
	
	var text string
	
	msg := BlankMessage(u.Update)
	chat_id := u.Update.Message.Chat.ID
	results, members, err := GetStats(chat_id)
	text_list := []string{"Rezultatele jocului pidoru zilei:"}
	
	if err == sql.ErrNoRows {
		msg.Text = "Nimei nu e inregistrat, statistica e goala"
		u.Bot.Send(msg)
		return
	}
	
	sort.Slice(members, func(i, j int) bool {
		return results[members[i].member_id].PidorCount < results[members[j].member_id].PidorCount
	})
	
	for _, memb := range members {
		txt := fmt.Sprintf("%v: %v", memb.full_name, results[memb.member_id].PidorCount)
		text_list = append(text_list, txt)
	}
	
	text = strings.Join(text_list, "\n")
	msg.Text = text
	
	u.Bot.Send(msg)
}