package main

import (
	"log"
	"github.com/The1Dani/faf_bot_go/internal"
)

func main() {
	// PostgreSQL //////////////////////////////
	db := internal.StartPostgreSQL()
	if db == nil {
		log.Println("Starting without PostgreSQL.")
	} else {
		defer db.Close()
		log.Println("The PostgreSQL has succesfully started!")
	}

	log.Println("Starting Telegram Bot")
	StartTelegramBot()

}
