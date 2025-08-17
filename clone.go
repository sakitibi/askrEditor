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

func cloneWiki(wikiSlug string) {
	url := "https://asakura-wiki.vercel.app/api/wiki/" + wikiSlug
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch wiki:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		fmt.Println("API error:", string(body))
		return
	}

	// 単一ページとしてアンマーシャル
	var page Page
	if err := json.Unmarshal(body, &page); err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}

	// ディレクトリ作成
	dir := page.WikiSlug
	os.MkdirAll(dir, 0755)

	// ファイル名は slug.askr
	filePath := filepath.Join(dir, page.Slug+".askr")
	if err := os.WriteFile(filePath, []byte(page.Content), 0644); err != nil {
		fmt.Println("Failed to write file:", err)
		return
	}

	fmt.Printf("✅ Saved %s\n", filePath)
}
