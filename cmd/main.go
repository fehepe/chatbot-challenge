package main

import (
	"fmt"
	"log"

	"github.com/fehepe/chatbot-challenge/pkg/db/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// dbConn := os.Getenv("CONN_DB")

	db, err := mysql.ConnectDB("root:root@tcp(db:3306)/chat_db")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":8080")
}
