package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func callAPI(method, wikiSlug, pageSlug string, body map[string]string, token string) (*http.Response, error) {
	url := "http://localhost:3000/api/wiki/" + wikiSlug + "/" + pageSlug

	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
