package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateKubebackupReq maps the request body for kubebackups
type CreateKubebackupReq struct {
	Name         string                    `json:"name"`
	Cluster      string                    `json:"cluster"`
	Policy       string                    `json:"policy,omitempty"`
	Pre_hooks    []string                  `json:"pre_hooks,omitempty"`
	Post_hooks   []string                  `json:"post_hooks,omitempty"`
	Trigger_type string                    `json:"trigger_type"`
	Source       CreateKubebackupReqSource `json:"source"`
}

// CreateKubebackupResp maps the POST response received from CloudCasa
type CreateKubebackupResp struct {
	Id           string   `json:"_id"`
	Name         string   `json:"name"`
	Cluster      string   `json:"cluster"`
	Policy       string   `json:"policy,omitempty"`
	Pre_hooks    []string `json:"pre_hooks"`
	Post_hooks   []string `json:"post_hooks"`
	Trigger_type string   `json:"trigger_type"`
	Updated      string   `json:"_updated"`
	Created      string   `json:"_created"`
	Etag         string   `json:"_etag"`
	Source       CreateKubebackupReqSource
	// Pause        bool     `json:"pause"`
}

// CreateKubebackupReqSource maps the 'source' dict for the request body
type CreateKubebackupReqSource struct {
	All_namespaces            bool     `json:"all_namespaces"`
	SnapshotPersistentVolumes bool     `json:"snapshotPersistentVolumes"`
	Namespaces                []string `json:"namespaces,omitempty"`
}

// GetKubebackupResp maps the GET response received from CloudCasa
type GetKubebackupResp struct {
	Id            string `json:"_id"`
	Name          string `json:"name"`
	Cc_user_email string `json:"cc_user_email"`
	Updated       string `json:"_updated"`
	Created       string `json:"_created"`
	Etag          string `json:"_etag"`
	Org_id        string `json:"org_id"`
}

// TODO: what do we need to return from run response?
type RunKubebackupResp struct {
	Id      string                  `json:"_id"`
	Cluster string                  `json:"cluster"`
	Name    string                  `json:"name"`
	Pause   bool                    `json:"pause"`
	Status  RunKubebackupRespStatus `json:"status"`
	Updated string                  `json:"_updated"`
	Created string                  `json:"_created"`
	Etag    string                  `json:"_etag"`
}

type RunKubebackupRespStatus struct {
	LastJobRunTime string              `json:"last_job_run_time"`
	Jobs           []map[string]string `json:"jobs"`
}

type JobStatus struct {
}

func (c *Client) RunKubebackup(backupId string, backupType string, retention int) (*RunKubebackupResp, error) {
	// Build request body
	reqBody := map[string]interface{}{
		"retention": map[string]int{
			"retainDays": retention,
		},
		"runBackup": true,
	}

	runReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// TODO: check if we need to add anything else for kubeoffload/kubebackup distinction
	runReq, err := http.NewRequest(http.MethodPost, c.ApiURL+backupType+"/"+backupId+"/action/run", bytes.NewBuffer(runReqBody))
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// POST to CloudCasa API
	runRespBody, err := c.doRequest(runReq)
	if err != nil {
		err = fmt.Errorf("error performing http request; %w", err)
		return nil, err
	}

	var runResp RunKubebackupResp
	if err := json.Unmarshal(runRespBody, &runResp); err != nil {
		return nil, err
	}

	return &runResp, nil
}

// TODO: move to job_client.go
func (c *Client) WatchJob(jobId string) (*JobStatus, error) {

	return nil, nil
}

// CreateKubebackup creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateKubebackup(reqBody CreateKubebackupReq) (*CreateKubebackupResp, error) {

	// Create rest request struct
	createReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	createReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"kubebackups", bytes.NewBuffer(createReqBody))
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

	var createResp CreateKubebackupResp
	if err := json.Unmarshal(createRespBody, &createResp); err != nil {
		return nil, err
	}

	return &createResp, nil
}

// GetKubebackup gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetKubebackup(kubebackupId string) (*GetKubebackupResp, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"kubebackups/"+kubebackupId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getKubebackupResp GetKubebackupResp
	if err := json.Unmarshal(getRespBody, &getKubebackupResp); err != nil {
		return nil, err
	}

	return &getKubebackupResp, nil
}

func (c *Client) UpdateKubebackup(kubebackupId string, reqBody interface{}, etag string) (*GetKubebackupResp, error) {
	// Create rest request struct
	putReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	putReq, err := http.NewRequest(http.MethodPut, c.ApiURL+"kubebackups/"+kubebackupId, bytes.NewBuffer(putReqBody))
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

	var putResp GetKubebackupResp
	if err := json.Unmarshal(putRespBody, &putResp); err != nil {
		return nil, err
	}

	return &putResp, nil
}

func (c *Client) DeleteKubebackup(kubebackupId string) error {
	delReq, err := http.NewRequest(http.MethodDelete, c.ApiURL+"kubebackups/"+kubebackupId, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(delReq)
	if err != nil {
		return err
	}

	return nil
}
