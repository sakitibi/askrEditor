package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
