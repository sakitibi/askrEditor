package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Page struct {
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

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

	// API は { pages: [...] } 形式を期待
	var result struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		// 単一ページの可能性もある
		var single Page
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			fmt.Println("Failed to parse JSON:", err)
			return
		}
		result.Pages = []Page{single}
	}

	if len(result.Pages) == 0 {
		fmt.Println("No pages found.")
		return
	}

	// wikiSlug ディレクトリを作成
	if err := os.MkdirAll(wikiSlug, 0755); err != nil {
		fmt.Println("Failed to create wiki directory:", err)
		return
	}

	// ページごとにファイルを作成
	for _, page := range result.Pages {
		filename := filepath.Join(wikiSlug, page.Slug+".askr")
		if err := os.WriteFile(filename, []byte(page.Content), 0644); err != nil {
			fmt.Printf("Failed to write %s: %v\n", filename, err)
		} else {
			fmt.Printf("Saved %s\n", filename)
		}
	}

	fmt.Printf("✅ Cloned wiki '%s' (%d pages)\n", wikiSlug, len(result.Pages))
}
