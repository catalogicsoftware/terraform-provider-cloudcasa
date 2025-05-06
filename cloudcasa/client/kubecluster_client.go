// Copyright 2025 Catalogic Software, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudcasa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// CreateKubeclusterResp maps the POST response received from CloudCasa
// We use different structs for create/get because 'status' field uses string
// for Create but struct for Get
type CreateKubeclusterResp struct {
	Id      string   `json:"_id"`
	Name    string   `json:"name"`
	Updated string   `json:"_updated"`
	Created string   `json:"_created"`
	Etag    string   `json:"_etag"`
	Status  string   `json:"_status"`
	Links   struct{} `json:"_links"`
}

type KubeclusterStatus struct {
	State     string `json:"state"`
	Agent_url string `json:"agentUrl"`
}

// GetKubeclusterResp maps the GET response received from CloudCasa
type GetKubeclusterResp struct {
	Id             string            `json:"_id"`
	Name           string            `json:"name"`
	Updated        string            `json:"_updated"`
	Created        string            `json:"_created"`
	Etag           string            `json:"_etag"`
	Status         KubeclusterStatus `json:"status"`
	Links          struct{}          `json:"_links"`
	BackupProvider struct {
		UserObjectstore string `json:"user_objectstore,omitempty"`
	} `json:"backup_provider,omitempty"`
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

func (c *Client) UpdateKubecluster(kubeclusterId string, reqBody interface{}, etag string) (*GetKubeclusterResp, error) {
	// Create rest request struct
	putReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	putReq, err := http.NewRequest(http.MethodPut, c.ApiURL+"kubeclusters/"+kubeclusterId, bytes.NewBuffer(putReqBody))
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

	var putResp GetKubeclusterResp
	if err := json.Unmarshal(putRespBody, &putResp); err != nil {
		return nil, err
	}

	return &putResp, nil
}

func (c *Client) DeleteKubecluster(kubeclusterId string) error {
	delReq, err := http.NewRequest(http.MethodDelete, c.ApiURL+"kubeclusters/"+kubeclusterId, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(delReq)
	if err != nil {
		return err
	}

	return nil
}

// Apply kubeagent spec to the kubecluster using kubectl and wait 5min for ACTIVE state
// Assumes Kubeconfig is set
// TODO: Validate kubeconfig
func (c *Client) ApplyKubeagent(clusterId string, agentUrl string) error {
	// Validate that KUBECONFIG is set
	if kconfig, exists := os.LookupEnv("KUBECONFIG"); exists != false {
		if kconfig == "" {
			return errors.New("KUBECONFIG environment variable is empty, but is required for auto_install. Set this variable to the path of a valid kubeconfig file.")
		}
	} else {
		return errors.New("KUBECONFIG environment variable is required for auto_install. Set this variable to the path of a valid kubeconfig file.")
	}

	// Check if the HTTP client is configured to skip TLS verification
	skipTLSVerify := false
	if transport, ok := c.HTTPClient.Transport.(*http.Transport); ok {
		if transport.TLSClientConfig != nil && transport.TLSClientConfig.InsecureSkipVerify {
			skipTLSVerify = true
		}
	}

	var err error
	if skipTLSVerify {
		// Use curl with -k flag to handle insecure SSL connections and pipe to kubectl
		curlCmd := exec.Command("bash", "-c", "curl -k " + agentUrl + " | kubectl apply -f -")
		_, err = curlCmd.Output()
	} else {
		// Use standard kubectl apply
		kubectlCmd := exec.Command("kubectl", "apply", "-f", agentUrl)
		_, err = kubectlCmd.Output()
	}

	if err != nil {
		return err
	}

	// Now wait for cluster to be ACTIVE
	// Wait 5min?
	for i := 1; i < 60; i++ {
		getKubeclusterResp, err := c.GetKubecluster(clusterId)
		if err != nil {
			return fmt.Errorf("error checking Kubecluster status after applying Agent manifest; %w", err)
		}

		if getKubeclusterResp.Status.State == "ACTIVE" {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timed out waiting for cluster to reach ACTIVE state")
}
