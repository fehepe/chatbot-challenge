package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	stockURL := os.Getenv("STOCK_URL")

	parsedUrl, err := url.Parse(strings.Replace(stockURL, "STOCKCODE", data["code"], 1))
	if err != nil {
		return err
	}

	client := http.Client{}

	resp, err := client.Get(parsedUrl.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','

	var response models.Stock

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !strings.Contains(data[0], "Symbol") && len(data) == 8 {
			response = models.Stock{
				Symbol: data[0],
				Date:   data[1],
				Time:   data[2],
				Open:   data[3],
				High:   data[4],
				Low:    data[5],
				Close:  data[6],
				Volume: data[7],
			}
		}
	}

	if response.Open == "N/D" {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"response": fmt.Sprintf("Your code: %s is not found.", response.Symbol),
		})
	} else {
		return c.JSON(fiber.Map{
			"response": fmt.Sprintf("%s quote is $%s per share.", response.Symbol, response.Open),
		})
	}

}
