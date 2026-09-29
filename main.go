package main

import (
	"log"
	_ "event_explorer/routers"
	"github.com/joho/godotenv"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	beego.Run()
}

