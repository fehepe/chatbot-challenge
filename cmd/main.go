package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	dbConn := os.Getenv("CONN_DB")

}
