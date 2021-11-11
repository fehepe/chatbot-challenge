package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	cookie := c.Cookies("jwt")

	_, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
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

	var response models.StockResponse

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !strings.Contains(data[0], "Symbol") && len(data) == 8 {
			response = models.StockResponse{
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

// map[string]string{
// 	"name":  "Toby",
// 	"email": "Toby@example.com",
// }

func Post(params map[string]string, endpoint string) ([]byte, error) {
	apiURL := os.Getenv("API")
	//Encode the data
	postBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://"+apiURL+endpoint, "application/json", responseBody)
	//Handle Error
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
