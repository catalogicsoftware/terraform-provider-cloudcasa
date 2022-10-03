package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO: XML?
var XML = "application/xml"
var JSON = "application/json"

func makeHttpRequest(url string, method string, contentType string, requestBody []byte, authToken string) []byte {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Authorization", "Bearer " + authToken)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: time.Second * 1000}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data
}
