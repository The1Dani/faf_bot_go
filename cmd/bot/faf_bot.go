package main

import (
	"database/sql"
	"log"

	"github.com/The1Dani/faf_bot_go/cmd/bot/commands"
	"github.com/The1Dani/faf_bot_go/internal"
)

var DB *sql.DB

func main() {
	// PostgreSQL //////////////////////////////

	DB = internal.StartPostgreSQL()
	commands.DB = DB
	
	if DB == nil {
		log.Println("Starting without PostgreSQL.")
	} else {
		defer DB.Close()
		log.Println("The PostgreSQL has succesfully started!")
	}

	log.Println("Starting Telegram Bot")
	StartTelegramBot()

}
