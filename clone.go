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
	Slug     string `json:"slug"`
	WikiSlug string `json:"wiki_slug"`
	Content  string `json:"content"`
}

// 単一ページをファイルに保存
func savePage(page Page) error {
	filePath := filepath.Join(page.WikiSlug, page.Slug+".askr")
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(page.Content), 0644)
}

// wikiSlug 内のすべてのページを取得して保存
func cloneWiki(wikiSlug string) {
	url := fmt.Sprintf("https://asakura-wiki.vercel.app/api/wiki/%s", wikiSlug)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch wiki:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	// ページ配列としてアンマーシャル
	var pages []Page
	if err := json.NewDecoder(resp.Body).Decode(&pages); err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}

	if len(pages) == 0 {
		fmt.Println("No pages found.")
		return
	}

	for _, page := range pages {
		if err := savePage(page); err != nil {
			fmt.Println("Failed to save page:", page.Slug, err)
			continue
		}
		fmt.Printf("✅ Saved %s/%s.askr\n", page.WikiSlug, page.Slug)
	}

	fmt.Println("✅ Clone finished")
}
