package iexcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IexcloudHistoricalPriceEntry struct {
	Close                float64 `json:"close"`
	High                 float64 `json:"high"`
	Low                  float64 `json:"low"`
	Open                 float64 `json:"open"`
	Symbol               string  `json:"symbol"`
	Volume               int64   `json:"volume"`
	Key                  string  `json:"key"`
	SubKey               string  `json:"subkey"`
	Date                 string  `json:"date"`
	UpdateTimestamp      int64   `json:"updated"`
	ChangeOverTime       float64 `json:"changeOverTime"`
	MarketChangeOverTime float64 `json:"marketChangeOverTime"`

	UOpen   float64 `json:"uOpen"`
	UClose  float64 `json:"uClose"`
	ULow    float64 `json:"uLow"`
	UHigh   float64 `json:"uHigh"`
	UVolume int64   `json:"uVolume"`

	FOpen   float64 `json:"fOpen"`
	FClose  float64 `json:"fClose"`
	FLow    float64 `json:"fLow"`
	FHigh   float64 `json:"fHigh"`
	FVolume int64   `json:"fVolume"`

	Label         string  `json:"label"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

func (this *Iexcloud) GetStockHistorialPrices(ticker string) ([]IexcloudHistoricalPriceEntry, error) {
	timeseriesEndpointUrl := fmt.Sprintf("%s/v1/stock/%s/chart/%s?token=%s", iexcloudBaseUrl, ticker, this.dateRange, this.apiKey)
	resp, err := http.Get(timeseriesEndpointUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []IexcloudHistoricalPriceEntry
	json.Unmarshal(bytes, &result)
	return result, nil
}
