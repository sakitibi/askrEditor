package main

import (
	"fmt"
	"os"
	"path/filepath"

	supabase "github.com/supabase-community/supabase-go"
)

func cloneWiki(client *supabase.Client, wikiSlug string) error {
	// Step 1: wikis テーブルで存在確認
	var wikis []struct {
		Slug string `json:"slug"`
	}
	_, err := client.From("wikis").
		Select("slug", "", false).
		Eq("slug", wikiSlug).
		ExecuteTo(&wikis)
	if err != nil {
		return fmt.Errorf("failed to fetch wiki: %w", err)
	}
	if len(wikis) == 0 {
		return fmt.Errorf("wiki '%s' not found", wikiSlug)
	}

	// Step 2: wiki_pages から slug, content を取得
	var pages []struct {
		Slug    string `json:"slug"`
		Content string `json:"content"`
	}
	_, err = client.From("wiki_pages").
		Select("slug, content", "", false).
		Eq("wiki_slug", wikiSlug).
		ExecuteTo(&pages)
	if err != nil {
		return fmt.Errorf("failed to fetch wiki_pages: %w", err)
	}

	// Step 3: ディレクトリ作成
	err = os.MkdirAll(wikiSlug, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Step 4: 各ページを .askr ファイルに保存
	for _, page := range pages {
		filePath := filepath.Join(wikiSlug, page.Slug+".askr")
		err := os.WriteFile(filePath, []byte(page.Content), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
		fmt.Printf("📄 %s\n", filePath)
	}

	fmt.Printf("✅ Cloned wiki '%s' (%d pages)\n", wikiSlug, len(pages))
	return nil
}
