package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateKubeoffloadReq maps the request body for kubeoffloads
type CreateKubeoffloadReq struct {
	Name              string              `json:"name"`
	Cluster           string              `json:"cluster"`
	Policy            string              `json:"policy,omitempty"`
	Trigger_type      string              `json:"trigger_type"`
	Backupdef         string              `json:"backupdef"`
	Run_backup        bool                `json:"run_backup"`
	Delete_snapshots  bool                `json:"delete_snapshots,omitempty"`
	Skip_live_copy    bool                `json:"skip_live_copy,omitempty"`
	Snapshot_longhorn bool                `json:"snapshot_longhorn,omitempty"`
	Offload_provider  KubeoffloadProvider `json:"offload_provider,omitempty"`
	Options           KubeoffloadOptions  `json:"options,omitempty"`
}

// TODO: check custom bucket settings
type KubeoffloadProvider struct {
	Type   string `json:"type,omitempty"`
	Region string `json:"region,omitempty"`
}

// TODO check Options
type KubeoffloadOptions struct {
	Example string `json:"example,omitempty"`
}

// Kubeoffload maps the GET response received from CloudCasa
type Kubeoffload struct {
	Id               string            `json:"_id"`
	Name             string            `json:"name"`
	Cluster          string            `json:"cluster"`
	Policy           string            `json:"policy"`
	Trigger_type     string            `json:"trigger_type"`
	Delete_snapshots bool              `json:"delete_snapshots"`
	Run_backup       bool              `json:"run_backup"`
	Backupdef        string            `json:"backupdef"`
	Updated          string            `json:"_updated"`
	Created          string            `json:"_created"`
	Etag             string            `json:"_etag"`
	Status           KubeoffloadStatus `json:"status"`
}

type KubeoffloadStatus struct {
	LastJobRunTime int64  `json:"last_job_run_time"`
	JobID          string `json:"jobid"`
	State          string `json:"state"`
}

func (c *Client) RunKubeoffload(backupId string, retention int) (*Kubeoffload, error) {
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

	runReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"kubeoffloads/"+backupId+"/action/run", bytes.NewBuffer(runReqBody))
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

	var runResp Kubeoffload
	if err := json.Unmarshal(runRespBody, &runResp); err != nil {
		return nil, err
	}

	return &runResp, nil
}

// CreateKubeoffload creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateKubeoffload(reqBody CreateKubeoffloadReq) (*Kubeoffload, error) {
	// Create rest request struct
	createReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	createReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"kubeoffloads", bytes.NewBuffer(createReqBody))
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// POST to CloudCasa API
	createResp, err := c.doRequest(createReq)
	if err != nil {
		err = fmt.Errorf("error performing http request; %w", err)
		return nil, err
	}

	var createRespBody Kubeoffload
	if err := json.Unmarshal(createResp, &createRespBody); err != nil {
		return nil, err
	}

	return &createRespBody, nil
}
