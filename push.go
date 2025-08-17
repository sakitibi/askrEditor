package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// pushWiki uploads all .askr files under wikiSlug directory
func pushWiki(wikiSlug string) {
	// ローカル保存されたトークンを取得
	token, err := getToken()
	if err != nil || token == "" {
		fmt.Println("❌ Not logged in. Please run: askreditor login <email> <password>")
		return
	}

	err = filepath.Walk(wikiSlug, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".askr") {
			return nil
		}

		// wikiSlug ディレクトリ以下の相対パスで slug を作る
		relPath, _ := filepath.Rel(wikiSlug, path)
		slug := strings.TrimSuffix(relPath, ".askr")
		slug = filepath.ToSlash(slug) // Windows 対応

		// ファイル内容を読む
		contentBytes, _ := os.ReadFile(path)
		body := map[string]string{
			"title":   slug,
			"content": string(contentBytes),
		}

		// API 呼び出し
		resp, err := callAPI("PUT", wikiSlug, slug, body, token)
		if err != nil {
			fmt.Println("❌ Failed:", slug, err)
			return nil
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		fmt.Println("✅ Pushed:", slug, string(data))
		return nil
	})
	if err != nil {
		fmt.Println("Push walk error:", err)
	}
}
