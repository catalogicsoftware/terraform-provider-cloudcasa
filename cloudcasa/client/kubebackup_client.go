package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Kubebackup maps the GET response received from CloudCasa
type Kubebackup struct {
	Id           string           `json:"_id,omitempty"`
	Name         string           `json:"name"`
	Cluster      string           `json:"cluster"`
	Policy       string           `json:"policy,omitempty"`
	Pre_hooks    []KubebackupHook `json:"pre_hooks,omitempty"`
	Post_hooks   []KubebackupHook `json:"post_hooks,omitempty"`
	Trigger_type string           `json:"trigger_type,omitempty"`
	Copydef      string           `json:"copydef,omitempty"`
	Updated      string           `json:"_updated,omitempty"`
	Created      string           `json:"_created,omitempty"`
	Etag         string           `json:"_etag,omitempty"`
	Source       KubebackupSource `json:"source"`
	Status       KubebackupStatus `json:"status,omitempty"`
}

type KubebackupStatus struct {
	LastJobRunTime int64               `json:"last_job_run_time,omitempty"`
	Jobs           []map[string]string `json:"jobs,omitempty"`
}

// KubebackupSource maps the 'source' dict for the request body
type KubebackupSource struct {
	All_namespaces            bool     `json:"all_namespaces"`
	SnapshotPersistentVolumes bool     `json:"snapshotPersistentVolumes"`
	Namespaces                []string `json:"namespaces,omitempty"`
}

type KubebackupHook struct {
	Template   bool     `json:"template"`
	Namespaces []string `json:"namespaces"`
	Hooks      []string `json:"hooks"`
}

// RunKubebackup runs the selected job using CloudCasa action/run API
func (c *Client) RunKubebackup(backupId string, retention int) (*Kubebackup, error) {
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

	runReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"kubebackups/"+backupId+"/action/run", bytes.NewBuffer(runReqBody))
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

	var runResp Kubebackup
	if err := json.Unmarshal(runRespBody, &runResp); err != nil {
		return nil, err
	}

	return &runResp, nil
}

// CreateKubebackup creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateKubebackup(reqBody Kubebackup) (*Kubebackup, error) {
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
	createResp, err := c.doRequest(createReq)
	if err != nil {
		err = fmt.Errorf("error performing http request; %w", err)
		return nil, err
	}

	var createRespBody Kubebackup
	if err := json.Unmarshal(createResp, &createRespBody); err != nil {
		return nil, err
	}

	return &createRespBody, nil
}

// GetKubebackup gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetKubebackup(kubebackupId string) (*Kubebackup, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"kubebackups/"+kubebackupId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getKubebackupResp Kubebackup
	if err := json.Unmarshal(getRespBody, &getKubebackupResp); err != nil {
		return nil, err
	}

	return &getKubebackupResp, nil
}

func (c *Client) UpdateKubebackup(kubebackupId string, reqBody Kubebackup, etag string) (*Kubebackup, error) {
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

	var putResp Kubebackup
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
