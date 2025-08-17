package main

import (
	"encoding/json"
	"io"
	"net/http"
)

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
