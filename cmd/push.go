package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
)

// pushWiki uploads all .askr files under wikiSlug directory
func PushWiki(wikiSlug string) {
	accessToken, err := auth.GetToken()
	if err != nil {
		colors.RedPrint1("❌", err)
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
			colors.RedPrint("Failed:", slug, err)
			return nil
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			colors.GreenPrint("✅ Pushed:", slug, string(data))
		} else {
			colors.RedPrint("❌ Failed to push:", slug, string(data))
		}
		return nil
	})

	if err != nil {
		colors.RedPrint1("Push walk error:", err)
	}
}
