package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func callAPI(method, wikiSlug, pageSlug string, body map[string]string, accessToken string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", apiBaseURL, wikiSlug, pageSlug)

	var reqBody *bytes.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonData)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	// ←ここが重要
	req.Header.Set("X-CLI", "true")

	return http.DefaultClient.Do(req)
}
