package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiBaseURL = "https://asakura-wiki.vercel.app/api/wiki"

func main() {
	http.HandleFunc("/replace", replaceHandler) // PUT
	http.HandleFunc("/clone", cloneHandler)     // POST
	http.HandleFunc("/delete", deleteHandler)   // DELETE
	http.HandleFunc("/get", getHandler)         // GET
	http.HandleFunc("/version", versionHandler)

	fmt.Println("✅ askreditor API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func getHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WikiSlug string `json:"wiki_slug"`
		Slug     string `json:"slug"`
		Token    string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.WikiSlug == "" {
		http.Error(w, "wiki_slug is required", http.StatusBadRequest)
		return
	}

	url := apiBaseURL + "/" + req.WikiSlug
	if req.Slug != "" {
		url += "/" + req.Slug
	}

	resp, err := requestAPI("GET", url, nil, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WikiSlug string `json:"wiki_slug"`
		Slug     string `json:"slug"`
		Token    string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.WikiSlug == "" || req.Slug == "" {
		http.Error(w, "wiki_slug and slug are required", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/%s/%s", apiBaseURL, req.WikiSlug, req.Slug)
	resp, err := requestAPI("DELETE", url, nil, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
