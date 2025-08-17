package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WikiSlug string `json:"wiki_slug"`
		Token    string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := callAPI("POST", req.WikiSlug, "", nil, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
