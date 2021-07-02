package Ticker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//GetDaysData returns all the trades data for a specificed stock.
func GetDaysData(apiKey, stock string, date *time.Time, count int) (string, error) {

	var (
		ticker     Ticker
		timeOffset = date
		err        error
	)

	if count != 0 {
		return GetMeteredTicker(apiKey, stock, date, count)
	}
	for {
		var returnedTickerData Ticker
		returnedTickerData, err = callForData(apiKey, timeOffset, stock, 50000)
		if err != nil {
			return "", err
		}

		if len(returnedTickerData.Results) == 0 {
			if len(ticker.Results) == 0 {

				return "", fmt.Errorf("No results for this day")
			}
			b, err := json.Marshal(ticker)
			if err != nil {
				return "", err
			}
			return string(b), nil
		}

		if len(ticker.Results) == 0 {
			ticker = returnedTickerData
			timeOffset = findLatestOffsetTime(returnedTickerData.Results)
			continue
		}

		timeOffset = findLatestOffsetTime(returnedTickerData.Results)
		ticker.Results = append(ticker.Results, returnedTickerData.Results...)
		ticker.ResultCount = ticker.ResultCount + returnedTickerData.ResultCount
	}

	return "", fmt.Errorf("Something unexpected happened on API call")
}

//GetMeteredTicker get the stocks from a limited number data
func GetMeteredTicker(apiKey, stock string, date *time.Time, limit int) (string, error) {

	var (
		ticker     Ticker
		timeOffset = date
		err        error
	)
	returnedTickerData, err := callForData(apiKey, timeOffset, stock, limit)
	if err != nil {
		return "", err
	}

	if len(returnedTickerData.Results) == 0 {
		return "", fmt.Errorf("No results for this day")
	}
	ticker = returnedTickerData
	b, err := json.Marshal(ticker)
	if err != nil {
		return "", err
	}
	return string(b), nil

}

func callForData(apiKey string, timeOffset *time.Time, stock string, limit int) (Ticker, error) {

	tlimit := time.Date(timeOffset.Year(), timeOffset.Month(), timeOffset.Day()+1, 0, 0, 0, 0, time.UTC)

	URL := fmt.Sprintf("https://api.polygon.io/v2/ticks/stocks/trades/%s/%s?timestamp=%d&timestampLimit=%d&limit=%d&apiKey=%s", stock, timeOffset.Format("2006-01-02"), timeOffset.UnixNano(), tlimit.UnixNano(), limit, apiKey)

	response, err := http.Get(URL)

	if err != nil {
		return Ticker{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Ticker{}, err
	}

	ticker := Ticker{}
	err = json.Unmarshal(responseData, &ticker)
	if err != nil {
		return Ticker{}, err
	}

	return ticker, nil
}

func findLatestOffsetTime(results []Result) *time.Time {

	// Polygon.io documentation states: Using the timestamp of the last result as the offset will give you the next page of results.
	result := results[len(results)-1]
	current := time.Unix(0, result.SIP)
	// Adding one nano second, we are able to not get the last result again.
	u, _ := time.ParseDuration(".001µs")
	offset := current.Add(u)

	return &offset
}
