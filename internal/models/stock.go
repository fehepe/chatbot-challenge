package models

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Stock struct {
	Symbol string
	Date   string
	Time   string
	Open   string
	High   string
	Low    string
	Close  string
	Volume string
}

func GetStock(data map[string]string) (Stock, error) {
	stockURL := os.Getenv("STOCK_URL")

	parsedUrl, err := url.Parse(strings.Replace(stockURL, "STOCKCODE", data["code"], 1))
	if err != nil {
		return Stock{}, err
	}

	client := http.Client{}

	resp, err := client.Get(parsedUrl.String())
	if err != nil {
		return Stock{}, err
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ','

	var response Stock

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !strings.Contains(data[0], "Symbol") && len(data) == 8 {
			response = Stock{
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

	return response, err
}
