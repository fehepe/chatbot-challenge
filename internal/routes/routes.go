package routes

import (
	"log"
	"net/http"

	"github.com/fehepe/chatbot-challenge/internal/controllers"
	"github.com/fehepe/chatbot-challenge/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gorilla/context"
)

func ApiSetup() {
	models.AutoMigrate()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
	app.Post("/api/stock", controllers.GetStock)

	app.Listen(":3000")
}

func WebSetup() {
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/logout", controllers.LogOutHandler)
	http.HandleFunc("/loginauth", controllers.LoginAuthHandler)
	http.HandleFunc("/register", controllers.RegisterHandler)
	http.HandleFunc("/registerauth", controllers.RegisterAuthHandler)
	http.HandleFunc("/ws", controllers.ChatServer)
	http.HandleFunc("/", controllers.IndexHandler)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("../../web"))))

	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux)))
}
