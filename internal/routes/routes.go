package routes

import (
	"log"
	"net/http"

	"github.com/fehepe/chatbot-challenge/internal/controllers"
	"github.com/fehepe/chatbot-challenge/pkg/db/hub"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApiSetup() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/stock", controllers.GetStock)
	app.Listen(":8080")
}

func WebSetup() {
	hubConn := hub.NewHub()
	go hubConn.Run()
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(hubConn, w, r)
	})
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("../../web"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
