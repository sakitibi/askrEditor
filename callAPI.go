package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func callAPI(method, wikiSlug, pageSlug string, body any, token string) (*http.Response, error) {
	url := apiBaseURL + "/" + wikiSlug
	if pageSlug != "" {
		url += "/" + pageSlug
	}

	var buf io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		buf = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	return client.Do(req)
}
