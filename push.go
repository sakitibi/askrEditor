package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// pushWiki uploads all .askr files under wikiSlug directory
func pushWiki(wikiSlug string) {
	accessToken, err := getToken()
	if err != nil {
		color.New(color.FgRed, color.Bold).Println("❌", err)
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

		contentBytes, _ := os.ReadFile(path)
		body := map[string]string{
			"title":   slug,
			"content": string(contentBytes),
		}

		resp, err := callAPI("PUT", wikiSlug, slug, body, accessToken)
		if err != nil {
			color.New(color.FgRed, color.Bold).Println("Failed:", slug, err)
			return nil
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			color.New(color.FgGreen, color.Bold).Println("✅ Pushed:", slug, string(data))
		} else {
			color.New(color.FgRed, color.Bold).Println("❌ Failed to push:", slug, string(data))
		}

		return nil
	})

	if err != nil {
		fmt.Println("Push walk error:", err)
	}
}
