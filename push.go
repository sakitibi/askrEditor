package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// pushWiki uploads all .askr files under wikiSlug directory
func pushWiki(wikiSlug string) {
	root := filepath.Join(".", wikiSlug)

	// ディレクトリが存在するか確認
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Println("Error: directory does not exist:", root)
		return
	}

	// ディレクトリツリーを走査
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリはスキップ
		if d.IsDir() {
			return nil
		}

		// .askr ファイルだけ対象
		if strings.HasSuffix(d.Name(), ".askr") {
			rel, _ := filepath.Rel(root, path) // wikiSlug からの相対パス
			pageSlug := strings.TrimSuffix(rel, ".askr")

			// ファイル内容を読む
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Failed to read:", path, err)
				return nil
			}

			// API 呼び出し
			payload := map[string]string{"content": string(content)}
			callAPI("PUT", wikiSlug, pageSlug, payload, "true")
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking wiki directory:", err)
	}
}
