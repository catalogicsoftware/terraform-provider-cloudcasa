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
	"fmt"
	"net/http"
)

// Policy maps the CloudCasa body for Policies
type Policy struct {
	Id        string           `json:"_id,omitempty"`
	Name      string           `json:"name"`
	Timezone  string           `json:"timezone"`
	Schedules []PolicySchedule `json:"schedules"`
	Updated   string           `json:"_updated,omitempty"`
	Created   string           `json:"_created,omitempty"`
	Etag      string           `json:"_etag,omitempty"`
}

// PolicySchedule maps the Schedule objects for Policies
type PolicySchedule struct {
	RetainDays int64          `json:"retainDays"`
	Locked     bool           `json:"locked"`
	Schedule   ScheduleStruct `json:"schedule"`
}

type ScheduleStruct struct {
	CronSpec string `json:"cronSpec"`
}

// CreatePolicy creates a resource in CloudCasa and returns a struct with important fields
func (c *Client) CreatePolicy(reqBody Policy) (*Policy, error) {
	// Create rest request struct
	createReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	createReq, err := http.NewRequest(http.MethodPost, c.ApiURL+"policies", bytes.NewBuffer(createReqBody))
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

	var createRespBody Policy
	if err := json.Unmarshal(createResp, &createRespBody); err != nil {
		return nil, err
	}

	return &createRespBody, nil
}

// GetPolicy gets a resource in CloudCasa and returns a struct with important fields
func (c *Client) GetPolicy(policyId string) (*Policy, error) {
	getReq, err := http.NewRequest(http.MethodGet, c.ApiURL+"policies/"+policyId, nil)
	if err != nil {
		return nil, err
	}

	getRespBody, err := c.doRequest(getReq)
	if err != nil {
		return nil, err
	}

	var getPolicyResp Policy
	if err := json.Unmarshal(getRespBody, &getPolicyResp); err != nil {
		return nil, err
	}

	return &getPolicyResp, nil
}

func (c *Client) UpdatePolicy(policyId string, reqBody Policy, etag string) (*Policy, error) {
	// Create rest request struct
	putReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	putReq, err := http.NewRequest(http.MethodPut, c.ApiURL+"policies/"+policyId, bytes.NewBuffer(putReqBody))
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

	var putResp Policy
	if err := json.Unmarshal(putRespBody, &putResp); err != nil {
		return nil, err
	}

	return &putResp, nil
}

func (c *Client) DeletePolicy(policyId string) error {
	delReq, err := http.NewRequest(http.MethodDelete, c.ApiURL+"policies/"+policyId, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(delReq)
	if err != nil {
		return err
	}

	return nil
}
