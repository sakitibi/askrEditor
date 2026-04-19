package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
	"golang.org/x/text/unicode/norm"
)

type wikiIndexResponse struct {
	PageSlugs []string `json:"page_slugs"`
}

func PushWiki(wikiSlug string) {
	accessToken, err := auth.GetToken()
	if err != nil {
		colors.RedPrint("❌ Auth Error:", err)
		os.Exit(1)
		return
	}

	// 1. サーバー側の slug 一覧取得（差分チェック用）
	resp, err := callAPI("GET", wikiSlug, "", nil, accessToken)
	if err != nil {
		colors.RedPrint("Failed to fetch wiki index:", err)
		os.Exit(1)
		return
	}
	defer resp.Body.Close()

	var index wikiIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		colors.RedPrint("Failed to parse wiki index:", err)
		os.Exit(1)
		return
	}

	serverSlugs := map[string]bool{}
	for _, s := range index.PageSlugs {
		serverSlugs[s] = true
	}

	localSlugs := map[string]bool{}

	// 2. ローカルファイルの走査とアップロード
	err = filepath.Walk(wikiSlug, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリ自体や .askr 以外はスキップ
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".askr") {
			return nil
		}

		// 相対パスを取得 (例: wikiSlug/folder/sub.askr -> folder/sub.askr)
		relPath, err := filepath.Rel(wikiSlug, path)
		if err != nil {
			return err
		}

		// 拡張子を除去し、OS固有の区切り文字を "/" に統一 (ネスト対策)
		slug := strings.TrimSuffix(relPath, ".askr")
		slug = filepath.ToSlash(slug)

		// ★ ここで NFC に正規化
		slug = norm.NFC.String(slug)

		localSlugs[slug] = true

		// ファイル読み込み
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rawText := string(contentBytes)
		title := slug
		content := rawText

		// TITLE: 行の抽出ロジック
		if strings.HasPrefix(rawText, "TITLE:") {
			lines := strings.SplitN(rawText, "\n", 2)
			title = strings.TrimSpace(strings.TrimPrefix(lines[0], "TITLE:"))
			if len(lines) > 1 {
				// 2行目以降（コンテンツ本体）から先頭の空行を除去
				content = strings.TrimLeft(lines[1], "\r\n")
			} else {
				content = ""
			}
		}

		body := map[string]string{
			"slug":    slug,
			"title":   title,
			"content": content,
		}

		// 新規作成か更新かを判定
		method := "PUT"
		if !serverSlugs[slug] {
			method = "POST"
		}

		resp, err := callAPI(method, wikiSlug, slug, body, accessToken)
		if err != nil {
			return fmt.Errorf("API call failed for %s: %w", slug, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			data, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("server error (%d) for %s: %s", resp.StatusCode, slug, string(data))
		}

		colors.GreenPrint("✅ [%s] %s", method, slug)
		return nil
	})

	if err != nil {
		colors.RedPrint("❌ Push failed:", err)
		os.Exit(1)
		return
	}

	// 3. 削除処理（ローカルに存在しないがサーバーにあるものを消す）
	for slug := range serverSlugs {
		// FrontPage は保護
		if slug == "FrontPage" || localSlugs[slug] {
			continue
		}

		resp, err := callAPI("DELETE", wikiSlug, slug, nil, accessToken)
		if err != nil {
			colors.RedPrint("⚠️ Delete failed:", slug, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == 200 {
			colors.GreenPrint("🗑 Deleted: %s", slug)
		}
	}
}
