package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

type CardInfoResponse struct {
	Token string `json:"token"`
	Phone string `json:"phone"`
}
type CardResponse struct {
	Card CardInfoResponse `json:"card"`
}
type CardTokenResponse struct {
	Error *string `json:"error"`
	Result *CardResponse `json:"result"`
}
func getHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+ os.Getenv("PAYMENT_AUTH"))
	return req
}
func GetCardToken(cardNumber, cardExpire string) (CardResponse, error) {
	data := map[string]interface{}{
		"method": "cards.get_phone",
		"params": map[string]interface{}{
			"number": cardNumber,
			"expire": cardExpire,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return CardResponse{}, err
	}
	req, err := http.NewRequest("POST", viper.GetString("payment.endpoint"),  bytes.NewBuffer(jsonData))
	if err != nil {
		return CardResponse{}, err
	}

	req = getHeader(req)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return CardResponse{}, errors.New("couldn't connect to broker api")
	}
	if resp.StatusCode != http.StatusOK {
		return CardResponse{}, errors.New("wrong request body")
	}
	defer resp.Body.Close()
	var responseData CardTokenResponse
	response, err := io.ReadAll(resp.Body)
	json.Unmarshal(response, &responseData)
	if responseData.Error != nil {
		return CardResponse{}, errors.New(*responseData.Error)
	}
	return *responseData.Result, err
}
