package controllers

import (
	"fmt"

	"github.com/fehepe/chatbot-challenge/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetStock(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "missing body for parameter",
		})
	}

	if data["code"] == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "missing code parameter",
		})
	}

	stockResponse, err := models.GetStock(data)
	if err != nil {
		return err
	}

	if stockResponse.Open == "N/D" {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"response": fmt.Sprintf("Your code: %s is not found.", stockResponse.Symbol),
		})
	} else {
		return c.JSON(fiber.Map{
			"response": fmt.Sprintf("%s quote is $%s per share.", stockResponse.Symbol, stockResponse.Open),
		})
	}

}
