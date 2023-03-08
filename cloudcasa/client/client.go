package cloudcasa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const ApiURL string = "https://api.staging.cloudcasa.io/api/v1/"
const JSON string = "application/json"

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
		return &c, nil
	}

	c.Apikey = *apikey

	// TODO: validate login with a test casa command
	// ar, err := c.SignIn()
	// if err != nil {
	// 	return nil, err
	// }

	return &c, nil
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

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
