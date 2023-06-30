package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
)

var (
	_client *http.Client
	baseURL string
	apiKey  string
)

type APIError struct {
	ResponseCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API returned code: %d", e.ResponseCode)
}

type HttpResponseBody struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

func GetApiClient(prefs fyne.Preferences) *http.Client {
	if _client == nil {
		_client = http.DefaultClient
	}

	baseURL = prefs.StringWithFallback("ssmurl", "https://ssmcloud.hostxtra.co.uk")
	apiKey = prefs.String("ssmapikey")

	return _client
}

func SendGetRequest(prefs fyne.Preferences, endpoint string, returnModel interface{}) error {

	GetApiClient(prefs)

	url := baseURL + endpoint

	fmt.Printf("#### GET #### url: %s\r\n", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-ssm-key", apiKey)

	r, err := _client.Do(req)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return &APIError{ResponseCode: r.StatusCode}
	}
	defer r.Body.Close()

	responseObject := HttpResponseBody{}

	json.NewDecoder(r.Body).Decode(&responseObject)

	if !responseObject.Success {
		return errors.New("api returned an error: " + responseObject.Error)
	}

	b, _ := json.Marshal(responseObject.Data)
	json.Unmarshal(b, returnModel)

	return nil
}

func SendPostRequest(prefs fyne.Preferences, endpoint string, bodyModel interface{}, returnModel interface{}) error {

	GetApiClient(prefs)

	bodyJSON, err := json.Marshal(bodyModel)

	if err != nil {
		return err
	}

	url := baseURL + endpoint

	fmt.Printf("#### POST #### url: %s, data: %s\r\n", url, bytes.NewBuffer(bodyJSON))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	req.Header.Set("x-ssm-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	r, err := _client.Do(req)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return &APIError{ResponseCode: r.StatusCode}
	}

	defer r.Body.Close()

	responseObject := HttpResponseBody{}

	json.NewDecoder(r.Body).Decode(&responseObject)

	if !responseObject.Success {
		fmt.Println(r.Body)
		return errors.New("api returned an error: " + responseObject.Error)
	}

	b, _ := json.Marshal(responseObject.Data)
	err = json.Unmarshal(b, returnModel)

	if err != nil {
		return err
	}

	return nil
}

func TestAPIConnection(prefs fyne.Preferences) error {
	var resModel interface{}
	return SendGetRequest(prefs, "/api/v1/account", &resModel)
}

func DownloadFile(prefs fyne.Preferences, url string, filePath string) error {
	GetApiClient(prefs)

	fmt.Printf("#### DOWNLOAD #### url: %s\r\n", url)

	req, _ := http.NewRequest("GET", url, nil)

	r, err := _client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, r.Body)
	return err
}

func SendGetRequestURL(prefs fyne.Preferences, url string, returnModel interface{}) error {
	GetApiClient(prefs)

	fmt.Printf("#### GET #### url: %s\r\n", url)

	req, _ := http.NewRequest("GET", url, nil)

	r, err := _client.Do(req)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return &APIError{ResponseCode: r.StatusCode}
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&returnModel)

	return nil
}
