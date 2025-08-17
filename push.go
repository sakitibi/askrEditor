package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// pushWiki uploads all .askr files under wikiSlug directory
func pushWiki(wikiSlug string) {
	err := filepath.Walk(wikiSlug, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // ディレクトリは無視
		}

		if !strings.HasSuffix(info.Name(), ".askr") {
			return nil // .askr 以外は無視
		}

		// wikiSlug/slug.askr の形に分解
		slug := strings.TrimSuffix(info.Name(), ".askr")
		contentBytes, _ := os.ReadFile(path)
		body := map[string]string{
			"title":   slug,
			"content": string(contentBytes),
		}

		resp, err := callAPI("PUT", wikiSlug, slug, body, "true") // X-CLI=true
		if err != nil {
			fmt.Println("Failed:", slug, err)
			return nil
		}
		defer resp.Body.Close()
		fmt.Println("✅ Pushed:", slug)
		return nil
	})
	if err != nil {
		fmt.Println("Push error:", err)
	}
}
