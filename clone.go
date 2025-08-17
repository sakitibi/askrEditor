package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	body, _ := io.ReadAll(resp.Body)

	// ページ配列としてパース
	type Page struct {
		Slug    string `json:"slug"`
		Content string `json:"content"`
	}
	var pages []Page
	if err := json.Unmarshal(body, &pages); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 各ページをファイルに書き出す
	for _, page := range pages {
		filename := fmt.Sprintf("%s.json", page.Slug)
		if err := os.WriteFile(filename, []byte(page.Content), 0644); err != nil {
			fmt.Printf("Failed to write %s: %v\n", filename, err)
		} else {
			fmt.Printf("Saved %s\n", filename)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(fmt.Sprintf("✅ Cloned wiki '%s' (%d pages)", req.WikiSlug, len(pages))))
}

// cloneWiki は wikiSlug 配下のページを取得してファイル出力
// cloneWiki は wikiSlug 配下のページを取得してファイル出力
func cloneWiki(wikiSlug string) {
	url := apiBaseURL + "/" + wikiSlug
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to call API:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ページ配列としてパース
	type Page struct {
		Slug    string `json:"slug"`
		Content string `json:"content"`
	}

	// API は { pages: [...] } の形式で返すと仮定
	var result struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		// 単一ページの可能性もあるので fallback
		var single Page
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			fmt.Println("Failed to parse JSON:", err)
			return
		}
		result.Pages = []Page{single}
	}

	// 各ページをファイルに書き出す
	for _, page := range result.Pages {
		filename := fmt.Sprintf("%s.json", page.Slug)
		if err := os.WriteFile(filename, []byte(page.Content), 0644); err != nil {
			fmt.Printf("Failed to write %s: %v\n", filename, err)
		} else {
			fmt.Printf("Saved %s\n", filename)
		}
	}

	fmt.Printf("✅ Cloned wiki '%s' (%d pages)\n", wikiSlug, len(result.Pages))
}
