package cloudcasa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// CreateKubeclusterResp maps the POST response received from CloudCasa
type CreateKubeclusterResp struct {
	Id            string   `json:"_id"`
	Name          string   `json:"name"`
	Cc_user_email string   `json:"cc_user_email"`
	Updated       string   `json:"_updated"`
	Created       string   `json:"_created"`
	Etag          string   `json:"_etag"`
	Org_id        string   `json:"org_id"`
	Status        string   `json:"_status"`
	Links         struct{} `json:"_links"`
}

type KubeclusterStatus struct {
	State     string `json:"state"`
	Agent_url string `json:"agentUrl"`
}

// GetKubeclusterResp maps the GET response received from CloudCasa
type GetKubeclusterResp struct {
	Id            string            `json:"_id"`
	Name          string            `json:"name"`
	Cc_user_email string            `json:"cc_user_email"`
	Updated       string            `json:"_updated"`
	Created       string            `json:"_created"`
	Etag          string            `json:"_etag"`
	Org_id        string            `json:"org_id"`
	Status        KubeclusterStatus `json:"status"`
	Links         struct{}          `json:"_links"`
}

// CreateKubecluster creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateKubecluster(reqBody interface{}) (*CreateKubeclusterResp, error) {

	// Create rest request struct
	createReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	createReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"kubeclusters", bytes.NewBuffer(createReqBody))
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// POST to CloudCasa API
	createRespBody, err := c.doRequest(createReq)
	if err != nil {
		err = fmt.Errorf("error performing http request; %w", err)
		return nil, err
	}

	var createResp CreateKubeclusterResp
	if err := json.Unmarshal(createRespBody, &createResp); err != nil {
		return nil, err
	}

	// Check that cluster resource was created in CloudCasa
	// TODO: Better failure check
	status := createResp.Status
	if status != "OK" {
		return nil, errors.New("received status NOT OK from CloudCasa")
	}

	return &createResp, nil
}

// GetKubecluster gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetKubecluster(kubeclusterId string) (*GetKubeclusterResp, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"kubeclusters/"+kubeclusterId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getKubeclusterResp GetKubeclusterResp
	if err := json.Unmarshal(getRespBody, &getKubeclusterResp); err != nil {
		return nil, err
	}

	return &getKubeclusterResp, nil
}
