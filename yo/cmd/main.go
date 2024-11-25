package main

import (
	"github.com/mymmrac/telego"
	"log"
	"myapp/environment"
	"myapp/pkg/database"
	processor "myapp/pkg/pocessor"
)

func main() {
	db, err := database.NewDataBase(environment.MustPath())
	if err != nil {
		log.Fatal("[ERR]from main:database didnt init -> %s", err)
	}
	//bot init
	bot, err := telego.NewBot(environment.MustToken())
	if err != nil {
		log.Fatalf("[ERR] problem with initialization -> %s", err)
	}

	log.Printf("[START] service is running")
	//geting updates
	for {
		updates, err := bot.UpdatesViaLongPolling(nil)
		if err != nil {
			log.Fatalf("[ERR] problem with getting updates -> %s ", err)
		}
		for update := range updates {
			processor.ListAndService(update, *bot, *db)
		}
	}
	// log.Printf("[RUN] service is runing")

}
