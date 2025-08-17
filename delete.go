package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WikiSlug string `json:"wiki_slug"`
		PageSlug string `json:"page_slug"`
		Token    string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.PageSlug == "" {
		req.PageSlug = "FrontPage"
	}

	resp, err := callAPI("DELETE", req.WikiSlug, req.PageSlug, nil, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
