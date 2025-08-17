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
		PageSlug string `json:"page_slug"`
		Content  string `json:"content"`
		Token    string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.PageSlug == "" {
		req.PageSlug = "FrontPage"
	}

	url := fmt.Sprintf("%s/%s/%s", apiBaseURL, req.WikiSlug, req.PageSlug)
	resp, err := requestAPI("PUT", url, map[string]string{"content": req.Content}, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
