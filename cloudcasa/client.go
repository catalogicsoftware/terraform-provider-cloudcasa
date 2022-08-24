package cloudcasa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const ApiURL string = "https://api.staging.cloudcasa.io"

type Client struct {
	ApiURL		string
	HTTPClient 	*http.Client
	Token		string
	Email		string
}

func NewClient(email, idToken *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ApiURL: ApiURL,
	}

	// If email/token are not provided, return empty client
	if email == nil || idToken == nil {
		return &c, nil
	}

	c.Token = *idToken
	c.Email = *email

	// TODO: validate login with a test casa command
	// ar, err := c.SignIn()
	// if err != nil {
	// 	return nil, err
	// }

	return &c, nil
}

func (c *Client) doRequest(req *http.Request, idToken *string) ([]byte, error) {
	// If no token is supplied in function call,
	// 		use token configured in Client object.
	token := c.Token

	if idToken != nil {
		token = *idToken
	}

	req.Header.Set("Authorization", token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
