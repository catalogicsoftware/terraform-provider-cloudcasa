package cloudcasa

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const ApiURL string = "https://home.cloudcasa.io/api/v1/"
const JSON string = "application/json"

const LogRequests bool = false

type Client struct {
	ApiURL     string
	HTTPClient *http.Client
	Apikey     string
}

func NewClient(apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiURL:     ApiURL,
	}

	// If apikey is not provided, return empty client
	if apikey == nil {
		return &c, errors.New("CloudCasa API key is missing")
	}

	c.Apikey = *apikey

	// TODO: validate login

	return &c, nil
}

// Log REST request/responses
func logRequest(body io.ReadCloser) error {
	if !LogRequests {
		return nil
	}

	parsedBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	//defer body.Close()

	err = ioutil.WriteFile("request_logs"+time.Now().String()+".txt", parsedBody, 0644)
	if err != nil {
		err = fmt.Errorf("CloudCasa Client: Error logging REST request: " + err.Error())
	}
	return err
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Assume client object has apikey
	apikey := c.Apikey

	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", JSON)
	req.Header.Set("Accept", JSON)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if LogRequests {
		err = ioutil.WriteFile("request_logs"+time.Now().String()+".txt", body, 0644)
		if err != nil {
			return nil, fmt.Errorf("CloudCasa Client: Error logging REST request: " + err.Error())
		}
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
