package cmd

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
)

// API から取得する slug 一覧用
type wikiIndexResponse struct {
	PageSlugs []string `json:"page_slugs"`
}

func PushWiki(wikiSlug string) {
	accessToken, err := auth.GetToken()
	if err != nil {
		colors.RedPrint("❌", err)
		return
	}

	// =========================
	// 1. サーバー側の slug 一覧取得
	// =========================
	resp, err := callAPI("GET", wikiSlug, "", nil, accessToken)
	if err != nil {
		colors.RedPrint("Failed to fetch wiki index:", err)
		return
	}
	defer resp.Body.Close()

	var index wikiIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		colors.RedPrint("Failed to parse wiki index:", err)
		return
	}

	serverSlugs := map[string]bool{}
	for _, s := range index.PageSlugs {
		serverSlugs[s] = true
	}

	localSlugs := map[string]bool{}

	// =========================
	// 2. ローカル → PUT / POST
	// =========================
	err = filepath.Walk(wikiSlug, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".askr") {
			return nil
		}

		relPath, _ := filepath.Rel(wikiSlug, path)
		slug := strings.TrimSuffix(relPath, ".askr")
		slug = filepath.ToSlash(slug)
		localSlugs[slug] = true

		contentBytes, _ := os.ReadFile(path)
		lines := strings.SplitN(string(contentBytes), "\n", 2)

		title := slug
		content := string(contentBytes)

		if strings.HasPrefix(lines[0], "TITLE:") {
			title = strings.TrimSpace(strings.TrimPrefix(lines[0], "TITLE:"))
			if len(lines) > 1 {
				content = strings.TrimLeft(lines[1], "\r\n")
			} else {
				content = ""
			}
			if title == "" {
				title = slug
			}
		}

		body := map[string]string{
			"title":   title,
			"content": content,
		}

		method := "PUT"
		if !serverSlugs[slug] {
			method = "POST"
		}

		resp, err := callAPI(method, wikiSlug, slug, body, accessToken)
		if err != nil {
			colors.RedPrint("❌ Failed:", slug, err)
			return nil
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			colors.GreenPrint("✅ %s: %s", method, slug)
		} else {
			colors.RedPrint("❌ Failed: %s\n%s", slug, string(data))
		}

		return nil
	})

	if err != nil {
		colors.RedPrint("Push walk error:", err)
		return
	}

	// =========================
	// 3. DELETE（ローカルに無いページ）
	// =========================
	for slug := range serverSlugs {
		if slug == "FrontPage" {
			continue
		}
		if localSlugs[slug] {
			continue
		}

		resp, err := callAPI("DELETE", wikiSlug, slug, nil, accessToken)
		if err != nil {
			colors.RedPrint("❌ Delete failed:", slug, err)
			continue
		}
		resp.Body.Close()

		colors.GreenPrint("🗑 Deleted: %s", slug)
	}
}
