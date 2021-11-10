package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fehepe/chatbot-challenge/internal/routes"
	"github.com/fehepe/chatbot-challenge/pkg/db/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	dbConn := os.Getenv("CONN_DB")

	db, err := mysql.ConnectDB(dbConn)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(db)

	routes.Setup(app)

	app.Listen(":8080")
}
