package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Kubeoffload maps the GET response received from CloudCasa
type Kubeoffload struct {
	Id                string              `json:"_id,omitempty"`
	Name              string              `json:"name"`
	Cluster           string              `json:"cluster"`
	Trigger_type      string              `json:"trigger_type"`
	Backupdef         string              `json:"backupdef"`
	Delete_snapshots  bool                `json:"delete_snapshots,omitempty"`
	Run_backup        bool                `json:"run_backup,omitempty"`
	Policy            string              `json:"policy,omitempty"`
	Skip_live_copy    bool                `json:"skip_live_copy,omitempty"`
	Snapshot_longhorn bool                `json:"snapshot_longhorn,omitempty"`
	Offload_provider  KubeoffloadProvider `json:"offload_provider,omitempty"`
	Options           KubeoffloadOptions  `json:"options,omitempty"`
	Updated           string              `json:"_updated,omitempty"`
	Created           string              `json:"_created,omitempty"`
	Etag              string              `json:"_etag,omitempty"`
	Status            KubeoffloadStatus   `json:"status,omitempty"`
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

type KubeoffloadStatus struct {
	LastJobRunTime int64  `json:"last_job_run_time,omitempty"`
	JobID          string `json:"jobid,omitempty"`
	State          string `json:"state,omitempty"`
}

// RunKubeoffload runs the selected job using CloudCasa action/run API
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
func (c *Client) CreateKubeoffload(reqBody Kubeoffload) (*Kubeoffload, error) {
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

// GetKubeoffload gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetKubeoffload(kubeoffloadId string) (*Kubeoffload, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"kubeoffloads/"+kubeoffloadId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getKubeoffloadResp Kubeoffload
	if err := json.Unmarshal(getRespBody, &getKubeoffloadResp); err != nil {
		return nil, err
	}

	return &getKubeoffloadResp, nil
}

func (c *Client) UpdateKubeoffload(kubeoffloadId string, reqBody Kubeoffload, etag string) (*Kubeoffload, error) {
	// Create rest request struct
	putReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	putReq, err := http.NewRequest(http.MethodPut, c.ApiURL+"kubeoffloads/"+kubeoffloadId, bytes.NewBuffer(putReqBody))
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

	var putResp Kubeoffload
	if err := json.Unmarshal(putRespBody, &putResp); err != nil {
		return nil, err
	}

	return &putResp, nil
}

func (c *Client) DeleteKubeoffload(kubeoffloadId string) error {
	delReq, err := http.NewRequest(http.MethodDelete, c.ApiURL+"kubeoffloads/"+kubeoffloadId, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(delReq)
	if err != nil {
		return err
	}

	return nil
}
