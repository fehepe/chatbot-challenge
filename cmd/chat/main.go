package main

import (
	"log"

	"github.com/fehepe/chatbot-challenge/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	routes.WebSetup()

}
