package cloudcasa

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const DefaultApiURL string = "https://home.cloudcasa.io/api/v1/"
const JSON string = "application/json"

const LogRequests bool = false

type Client struct {
	ApiURL     string
	HTTPClient *http.Client
	Apikey     string
}

func NewClient(apikey *string, cloudcasaUrl *string, allowInsecureTLS bool) (*Client, error) {
	// Create HTTP client with or without TLS verification
	httpClient := &http.Client{Timeout: 10 * time.Second}
	
	// If allowInsecureTLS is true, configure the client to skip certificate verification
	if allowInsecureTLS {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	c := Client{
		HTTPClient: httpClient,
	}

	// If apikey is not provided, return empty client
	if apikey == nil {
		return &c, errors.New("CloudCasa API key is missing")
	}

	c.Apikey = *apikey

	// Set API URL - use provided value or default
	if cloudcasaUrl != nil && *cloudcasaUrl != "" {
		c.ApiURL = *cloudcasaUrl
		// Ensure URL ends with /api/v1/
		if c.ApiURL[len(c.ApiURL)-1:] != "/" {
			c.ApiURL = c.ApiURL + "/"
		}
		if len(c.ApiURL) < 8 || c.ApiURL[len(c.ApiURL)-8:] != "/api/v1/" {
			c.ApiURL = c.ApiURL + "api/v1/"
		}
	} else {
		c.ApiURL = DefaultApiURL
	}

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
