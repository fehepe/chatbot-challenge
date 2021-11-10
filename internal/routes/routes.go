package routes

import (
	"github.com/fehepe/chatbot-challenge/internal/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/", controller.Hello)

}
