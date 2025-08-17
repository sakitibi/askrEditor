package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const apiBaseURL = "https://asakura-wiki.vercel.app/api/wiki"

func main() {
	http.HandleFunc("/replace", replaceHandler) // PUT
	http.HandleFunc("/clone", cloneHandler)     // POST
	http.HandleFunc("/delete", deleteHandler)   // DELETE
	http.HandleFunc("/get", getHandler)         // GET
	http.HandleFunc("/version", versionHandler)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version))
}

func requestAPI(method, url string, body any, token string) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
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
