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
		colors.RedPrint("❌", err)
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

		relPath, _ := filepath.Rel(wikiSlug, path)
		slug := strings.TrimSuffix(relPath, ".askr")
		slug = filepath.ToSlash(slug)

		contentBytes, _ := os.ReadFile(path)
		lines := strings.SplitN(string(contentBytes), "\n", 2)

		var title string
		var content string

		if len(lines) > 0 && strings.HasPrefix(lines[0], "TITLE:") {
			title = strings.TrimSpace(strings.TrimPrefix(lines[0], "TITLE:"))
			if len(lines) > 1 {
				content = lines[1]
			} else {
				content = ""
			}
		} else {
			title = slug
			content = string(contentBytes)
		}

		body := map[string]string{
			"title":   title,
			"content": content,
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
		colors.RedPrint("Push walk error:", err)
	}
}
