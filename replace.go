package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func replaceHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WikiSlug string `json:"wiki_slug"`
		Slug     string `json:"slug"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Token    string `json:"token"` // Bearer トークン
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
	body := map[string]string{
		"title":   req.Title,
		"content": req.Content,
	}

	resp, err := requestAPI("PUT", url, body, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
