package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ======================
// pushWiki はローカルファイルを wiki に同期
// ======================
func pushWiki(wikiSlug string) {
	baseDir := wikiSlug
	filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Walk error:", err)
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".askr") {
			return nil
		}
		rel, _ := filepath.Rel(baseDir, path)
		pageSlug := strings.TrimSuffix(rel, ".askr")
		contentBytes, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("Failed to read file:", path, err)
			return nil
		}
		content := string(contentBytes)
		callAPI("PUT", wikiSlug, pageSlug, map[string]string{"content": content}, "true")
		fmt.Printf("✅ Pushed %s/%s\n", wikiSlug, pageSlug)
		return nil
	})
}
