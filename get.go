package main

import (
	"fmt"
	"io"
	"net/http"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	wikiSlug := r.URL.Query().Get("wiki_slug")
	pageSlug := r.URL.Query().Get("page_slug")
	if pageSlug == "" {
		pageSlug = "FrontPage"
	}

	url := fmt.Sprintf("%s/%s/%s", apiBaseURL, wikiSlug, pageSlug)
	resp, err := requestAPI("GET", url, nil, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
