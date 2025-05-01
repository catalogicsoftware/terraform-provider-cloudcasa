package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Objectstore maps the CloudCasa body for Objectstores
type Objectstore struct {
	Id                  string       `json:"_id,omitempty"`
	Name                string       `json:"name"`
	Private             bool         `json:"private"`
	ProxyClusterList    []string     `json:"proxy_cluster_list,omitempty"`
	ProviderType        string       `json:"provider_type"`
	BucketName          string       `json:"bucket_name,omitempty"`
	Region              string       `json:"region,omitempty"`
	SkipTlsCertificateValidation bool `json:"skip_tls_certificate_validation,omitempty"`
	S3Provider          S3Provider   `json:"s3provider,omitempty"`
	Updated             string       `json:"_updated,omitempty"`
	Created             string       `json:"_created,omitempty"`
	Etag                string       `json:"_etag,omitempty"`
}

// S3Provider maps the S3 provider dict in the request body 
type S3Provider struct {
	Endpoint			string					`json:"endpoint,omitempty"`
	Credentials			S3ProviderCredentials	`json:"credentials,omitempty"`
	Cloud				string					`json:"cloud,omitempty"`
	ResourceGroupName	string					`json:"resource_group_name,omitempty"`
	StorageAccountName	string					`json:"storage_account_name,omitempty"`	
}

// S3ProviderCredentials maps the credentials dict in the S3Provider struct
type S3ProviderCredentials struct {
	SubscriptionId string `json:"subscription_id,omitempty"`
	TenantId       string `json:"tenant_id,omitempty"`
	ClientId       string `json:"client_id,omitempty"`
	ClientSecret   string `json:"client_secret,omitempty"`
	AccessKey      string `json:"access_key,omitempty"`
	SecretKey      string `json:"secret_key,omitempty"`
}

// CreateObjectstore creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateObjectstore(reqBody Objectstore) (*Objectstore, error) {
	// Create rest request struct
	createReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	createReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"objectstores", bytes.NewBuffer(createReqBody))
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

	var createResp Objectstore
	if err := json.Unmarshal(createRespBody, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

// GetObjectstore gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetObjectstore(objectstoreId string) (*Objectstore, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"objectstores/"+objectstoreId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getObjectstoreResp Objectstore
	if err := json.Unmarshal(getRespBody, &getObjectstoreResp); err != nil {
		return nil, err
	}

	return &getObjectstoreResp, nil
}

// UpdateObjectstore updates a resource in CloudCasa and returns a struct with important fields
func (c *Client) UpdateObjectstore(objectstoreId string, reqBody Objectstore, etag string) (*Objectstore, error) {
	// Create rest request struct
	putReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	putReq, err := http.NewRequest(http.MethodPut, c.ApiURL+"objectstores/"+objectstoreId, bytes.NewBuffer(putReqBody))
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// PUT requests require matching etag
	putReq.Header.Set("If-Match", etag)

	// PUT to CloudCasa API
	putRespBody, err := c.doRequest(putReq)
	if err != nil {
		err = fmt.Errorf("error performing http request; %w", err)
		return nil, err
	}

	var putResp Objectstore
	if err := json.Unmarshal(putRespBody, &putResp); err != nil {
		return nil, err
	}

	return &putResp, nil
}

// DeleteObjectstore deletes a resource in CloudCasa
func (c *Client) DeleteObjectstore(objectstoreId string) error {
	delReq, err := http.NewRequest(http.MethodDelete, c.ApiURL+"objectstores/"+objectstoreId, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(delReq)
	if err != nil {
		return err
	}

	return nil
} 