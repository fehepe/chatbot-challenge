package stock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fehepe/chatbot-challenge/internal/models"
)

func GetStockFromAPI(code string) (models.StockResponse, error) {
	params := map[string]string{
		"code": code,
	}
	response, err := Post(params, "/api/stock")
	if err != nil {
		fmt.Println("error doing the login request")
		return models.StockResponse{}, err
	}

	var resp models.StockResponse
	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
		return models.StockResponse{}, err
	}

	return resp, nil
}

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
