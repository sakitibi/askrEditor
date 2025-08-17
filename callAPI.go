package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func callAPI(method, wikiSlug, pageSlug string, body map[string]string, token string) (*http.Response, error) {
	var url string
	if pageSlug != "" {
		url = fmt.Sprintf("%s/%s/%s", apiBaseURL, wikiSlug, pageSlug)
	} else {
		url = fmt.Sprintf("%s/%s/FrontPage", apiBaseURL, wikiSlug)
	}

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
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return http.DefaultClient.Do(req)
}
