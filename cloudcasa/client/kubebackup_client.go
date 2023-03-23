package cloudcasa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TODO: should we use omitempty for json fields?

// CreateKubebackupReq maps the request body for kubebackups
type CreateKubebackupReq struct {
	Name         string           `json:"name"`
	Cluster      string           `json:"cluster"`
	Policy       string           `json:"policy"`
	Pre_hooks    []KubebackupHook `json:"pre_hooks"`
	Post_hooks   []KubebackupHook `json:"post_hooks"`
	Trigger_type string           `json:"trigger_type"`
	Source       KubebackupSource `json:"source"`
}

type KubebackupHook struct {
	Template   bool     `json:"template"`
	Namespaces []string `json:"namespaces"`
	Hooks      []string `json:"hooks"`
}

// GetKubebackupResp maps the GET response received from CloudCasa
type GetKubebackupResp struct {
	Id           string           `json:"_id"`
	Name         string           `json:"name"`
	Cluster      string           `json:"cluster"`
	Policy       string           `json:"policy"`
	Pre_hooks    []KubebackupHook `json:"pre_hooks"`
	Post_hooks   []KubebackupHook `json:"post_hooks"`
	Trigger_type string           `json:"trigger_type"`
	Updated      string           `json:"_updated"`
	Created      string           `json:"_created"`
	Etag         string           `json:"_etag"`
	Source       KubebackupSource `json:"source"`
	Status       KubebackupStatus `json:"status"`
	Org_id       string           `json:"org_id"`
}

type KubebackupStatus struct {
	LastJobRunTime int64               `json:"last_job_run_time"`
	Jobs           []map[string]string `json:"jobs"`
}

// KubebackupSource maps the 'source' dict for the request body
type KubebackupSource struct {
	All_namespaces            bool     `json:"all_namespaces"`
	SnapshotPersistentVolumes bool     `json:"snapshotPersistentVolumes"`
	Namespaces                []string `json:"namespaces"`
}
type GetJobsResp struct {
	Items []GetJobResp `json:"_items"`
}

type GetJobResp struct {
	Id         string   `json:"_id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	State      string   `json:"state"`
	Start_Time int64    `json:"start_time"`
	Cluster    string   `json:"cluster"`
	Jobrunner  string   `json:"jobrunner"`
	Activity   []string `json:"activity"`
	Updated    string   `json:"_updated"`
	Created    string   `json:"_created"`
	Etag       string   `json:"_etag"`
}

func (c *Client) RunKubebackup(backupId string, backupType string, retention int) (*GetKubebackupResp, error) {
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

	var runResp GetKubebackupResp
	if err := json.Unmarshal(runRespBody, &runResp); err != nil {
		return nil, err
	}

	return &runResp, nil
}

// TODO: move to job_client.go
func (c *Client) GetJobFromBackupdef(backupId string, lastRunTime int64) (*GetJobResp, error) {
	// Create HTTP Request
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"jobs", nil)
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// Set HTTP queries
	// TODO: returning 0 jobs
	// maybe because of last job run time?
	queries := getReq.URL.Query()
	whereQueryString := fmt.Sprintf("{\"type\":{\"$nin\":[\"DELETE_BACKUP\",\"AWSRDS_BACKUP_DELETE\",\"AGENT_UPDATE\"]},\"backupdef\":\"%s\",\"start_time\":{\"$gt\":%d}}", backupId, lastRunTime)
	//	whereQueryString := fmt.Sprintf("{\"type\":{\"$nin\":[\"DELETE_BACKUP\",\"AWSRDS_BACKUP_DELETE\",\"AGENT_UPDATE\"]},\"state\":{\"$in\":[\"RUNNING\",\"CATALOG\"]}}&sort=-start_time&backupdef_name=\"%s\"&\"status.last_job_run_time\":{\"$gt\":%s}", backupName, lastRunTime)
	queries.Add("where", whereQueryString)
	queries.Add("sort", "-start_time")
	queries.Add("page", "1")
	queries.Add("max_results", "5")
	getReq.URL.RawQuery = queries.Encode()
	//fmt.Println(getReq.URL.String())

	// TODO: do we need this?
	cookie := &http.Cookie{
		Name:  "auth0.is.authenticated",
		Value: "true",
	}
	getReq.AddCookie(cookie)

	foundJob := false
	var getJobsResp GetJobsResp

	// wait 1min for job to appear
	for i := 1; i < 12; i++ {
		// GET from CloudCasa API
		getJobsRespBody, err := c.doRequest(getReq)
		if err != nil {
			err = fmt.Errorf("error performing http request; %w", err)
			return nil, err
		}

		// Unmarshall response
		if err := json.Unmarshal(getJobsRespBody, &getJobsResp); err != nil {
			return nil, err
		}

		// If any jobs match the filter, we found the job
		if len(getJobsResp.Items) > 0 {
			foundJob = true
			break
		}

		time.Sleep(5 * time.Second)
	}

	if !foundJob {
		return nil, fmt.Errorf("could not find job created by kubebackup %s", backupId)
	}

	return &getJobsResp.Items[0], nil
}

func (c *Client) WatchJobUntilComplete(jobId string) (*GetJobResp, error) {
	// Create HTTP Request
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"jobs/"+jobId, nil)
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	doneStates := []string{"COMPLETED", "SKIPPED", "PARTIAL"}

	// Wait 5 minutes for job to complete? TODO: decide job watch timeout
	for i := 1; i < 60; i++ {
		getResp, err := c.doRequest(getReq)
		if err != nil {
			err = fmt.Errorf("error performing http request; %w", err)
			return nil, err
		}

		var jobRespBody GetJobResp
		if err := json.Unmarshal(getResp, &jobRespBody); err != nil {
			return nil, err
		}

		for _, v := range doneStates {
			// Job is in a completed state
			if v == jobRespBody.State {
				return &jobRespBody, nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("timed out waiting for job %s to complete", jobId)
}

// CreateKubebackup creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreateKubebackup(reqBody CreateKubebackupReq) (*GetKubebackupResp, error) {
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

	var createRespBody GetKubebackupResp
	if err := json.Unmarshal(createResp, &createRespBody); err != nil {
		return nil, err
	}

	return &createRespBody, nil
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
