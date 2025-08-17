package main

import (
	"io"
	"net/http"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	wikiSlug := r.URL.Query().Get("wiki_slug")
	pageSlug := r.URL.Query().Get("page_slug")
	if pageSlug == "" {
		pageSlug = "FrontPage"
	}

	resp, err := callAPI("GET", wikiSlug, pageSlug, nil, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
