package main

import (
	"log"
	"os"

	"github.com/fehepe/chatbot-challenge/internal/routes"
	"github.com/fehepe/chatbot-challenge/pkg/db/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	dbConn := os.Getenv("CONN_DB")

	if err := mysql.ConnectDB(dbConn); err != nil {
		log.Fatal(err.Error())
	}

	routes.ApiSetup()
}
