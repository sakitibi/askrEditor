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

// savePage は page を wikiSlug/slug.askr に保存
func savePage(page Page) error {
	filePath := filepath.Join(page.WikiSlug, page.Slug+".askr")
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(page.Content), 0644)
}

// fetchPage は API からページを取得
func fetchPage(wikiSlug, pageSlug string) (*Page, error) {
	url := fmt.Sprintf("https://asakura-wiki.vercel.app/api/wiki/%s", wikiSlug)
	if pageSlug != "" {
		url += "/" + pageSlug
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &page, nil
}

// fetchSlugs は API から wikiSlug のページ一覧を取得
func fetchSlugs(wikiSlug string) ([]string, error) {
	url := fmt.Sprintf("https://asakura-wiki.vercel.app/api/wiki/%s", wikiSlug)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slugs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// page_slugs を取り出す
	var data struct {
		WikiSlug  string   `json:"wiki_slug"`
		PageSlugs []string `json:"page_slugs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data.PageSlugs, nil
}

func cloneWiki() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: askreditor clone <wikiSlug>")
		return
	}
	wikiSlug := os.Args[2]

	// 1. slug 一覧を取得
	slugs, err := fetchSlugs(wikiSlug)
	if err != nil {
		fmt.Println("Failed to fetch slug list:", err)
		return
	}

	// 2. 各ページを保存
	for _, slug := range slugs {
		page, err := fetchPage(wikiSlug, slug)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := savePage(*page); err != nil {
			fmt.Println("Failed to save page:", err)
			continue
		}
		fmt.Printf("✅ Saved %s/%s.askr\n", page.WikiSlug, page.Slug)
	}
}
