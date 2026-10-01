package main

import (
	"log"

	"event_explorer/services"
	"event_explorer/utils"

	_ "event_explorer/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

var runServer = beego.Run

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	config := utils.LoadConfig()

	services.Initialize(
		config.GooglePlacesAPIKey,
		config.TicketmasterAPIKey,
	)

	runServer()
}
