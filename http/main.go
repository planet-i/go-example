package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type ApiCall struct {
	REQ_URL   string
	APP_KEY   string
	APP_TOKEN string
}

func NewApiCall() *ApiCall {
	return &ApiCall{
		REQ_URL:   "https://gwapi.mabangerp.com/api/v2",
		APP_KEY:   "200636",                           // Replace with your APP_KEY
		APP_TOKEN: "58ab5a9a70484ab9b0209b012636c384", // Replace with your APP_TOKEN
	}
}

func HMAC256(c, key string) string {
	sig := hmac.New(sha256.New, []byte(key))
	sig.Write([]byte(c))
	return hex.EncodeToString(sig.Sum(nil))
}

// API Request Call
func (api *ApiCall) Call(apiName string, reqParams map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"api":       apiName,
		"appkey":    api.APP_KEY,
		"data":      reqParams,
		"timestamp": int(time.Now().Unix()),
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %v", err)
	}

	// Generate Authorization header using HMAC-SHA256
	authorization := HMAC256(string(dataJSON), api.APP_TOKEN)

	// Prepare headers
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Requested-With": "XMLHttpRequest",
		"Authorization":    authorization,
	}

	// Use Resty client to send POST request
	client := resty.New()
	resp, err := client.R().
		SetHeaders(headers).
		SetBody(dataJSON).
		Post(api.REQ_URL)

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func main() {
	// Create new API client
	apicall := NewApiCall()

	// Make the API call
	result, err := apicall.Call("stock-do-search-sku-list-new", map[string]interface{}{})
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	// Print the result
	fmt.Println(result)
}
