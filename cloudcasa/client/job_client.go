package cloudcasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Jobs struct {
	Items []Job `json:"_items"`
}

type Job struct {
	Id         string   `json:"_id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	State      string   `json:"state"`
	Start_Time int64    `json:"start_time"`
	Cluster    string   `json:"cluster"`
	Jobrunner  string   `json:"jobrunner,omitempty"`
	Activity   []string `json:"activity,omitempty"`
	Updated    string   `json:"_updated"`
	Created    string   `json:"_created"`
	Etag       string   `json:"_etag"`
}

// GetJobFromBackupdef returns job objects for corresponding backupdef
func (c *Client) GetJobFromBackupdef(backupId string, lastRunTime int64) (*Job, error) {
	// Create HTTP Request
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"jobs", nil)
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	// Set HTTP queries
	queries := getReq.URL.Query()
	whereQueryString := fmt.Sprintf("{\"type\":{\"$nin\":[\"DELETE_BACKUP\",\"AWSRDS_BACKUP_DELETE\",\"AGENT_UPDATE\"]},\"backupdef\":\"%s\",\"start_time\":{\"$gt\":%d}}", backupId, lastRunTime)
	//	whereQueryString := fmt.Sprintf("{\"type\":{\"$nin\":[\"DELETE_BACKUP\",\"AWSRDS_BACKUP_DELETE\",\"AGENT_UPDATE\"]},\"state\":{\"$in\":[\"RUNNING\",\"CATALOG\"]}}&sort=-start_time&backupdef_name=\"%s\"&\"status.last_job_run_time\":{\"$gt\":%s}", backupName, lastRunTime)
	queries.Add("where", whereQueryString)
	queries.Add("sort", "-start_time")
	queries.Add("page", "1")
	queries.Add("max_results", "5")
	getReq.URL.RawQuery = queries.Encode()

	cookie := &http.Cookie{
		Name:  "auth0.is.authenticated",
		Value: "true",
	}
	getReq.AddCookie(cookie)

	foundJob := false
	var getJobsResp Jobs

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

// WatchJobUntilComplete finds the job with given ID and waits until job completes to return
func (c *Client) WatchJobUntilComplete(jobId string) (*Job, error) {
	// Create HTTP Request
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"jobs/"+jobId, nil)
	if err != nil {
		err = fmt.Errorf("error creating http request; %w", err)
		return nil, err
	}

	doneStates := []string{"COMPLETED", "SKIPPED", "PARTIAL", "CANCELLED"}

	// Wait 5 minutes for job to complete
	for i := 1; i < 60; i++ {
		getResp, err := c.doRequest(getReq)
		if err != nil {
			err = fmt.Errorf("error performing http request; %w", err)
			return nil, err
		}

		var jobRespBody Job
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

	return nil, fmt.Errorf("timed out waiting for job %s to complete. The terraform provider will only wait 5minutes for a running job. For more details and job logs, check CloudCasa UI", jobId)
}
