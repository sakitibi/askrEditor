package wiki

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
)

type Page struct {
	Slug     string `json:"slug"`
	WikiSlug string `json:"wiki_slug"`
	Title    string `json:"title"`
	Content  string `json:"content"` // API v2 では Base64(Gzipped) 文字列
}

// decodeContent は Base64(Gzip) 形式の文字列を元のテキストに復元します
func decodeContent(encodedStr string) (string, error) {
	// 1. Base64 デコード
	compressed, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	// 2. Gzip 解凍
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", fmt.Errorf("gzip reader error: %w", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read decompressed data: %w", err)
	}

	return string(decompressed), nil
}

func callAPIWikis(accessToken string) ([]string, error) {
	// エンドポイントを /api/wikis に維持 (Blob一覧取得用)
	apiURL := "https://asakura-wiki.vercel.app/api/wikis"
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result []string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func savePage(page Page) error {
	filePath := filepath.Join(page.WikiSlug, page.Slug+".askr")

	// ディレクトリを階層ごと作成
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data := fmt.Sprintf("TITLE:%s\n%s", page.Title, page.Content)
	return os.WriteFile(filePath, []byte(data), 0644)
}

func fetchPage(wikiSlug, pageSlug string) (*Page, error) {
	// API v2 の URL 組み立て
	url := fmt.Sprintf("%s/%s/%s", auth.ApiBaseURL, wikiSlug, pageSlug)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d", resp.StatusCode)
	}

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 重要: ここでバイナリをテキストにデコードする
	rawText, err := decodeContent(page.Content)
	if err != nil {
		return nil, fmt.Errorf("decode error for %s: %w", pageSlug, err)
	}
	page.Content = rawText
	page.WikiSlug = wikiSlug
	page.Slug = pageSlug

	return &page, nil
}

func fetchSlugs(wikiSlug string) ([]string, error) {
	// 叩いているURLを可視化する
	url := fmt.Sprintf("%s/%s", auth.ApiBaseURL, wikiSlug)
	fmt.Printf("DEBUG: Fetching slugs from: %s\n", url) // これで404のURLを特定

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// ボディを取得して詳細を確認する
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error 404 at %s: %s", url, string(body))
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d", resp.StatusCode)
	}

	var data struct {
		PageSlugs []string `json:"page_slugs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data.PageSlugs, nil
}

func CloneWiki(wikiSlug string) {
	slugs, err := fetchSlugs(wikiSlug)
	if err != nil {
		colors.RedPrint("Failed to fetch slug list: %s", err)
		return
	}
	if len(slugs) == 0 {
		colors.RedPrint("%s has no pages", wikiSlug)
		return
	}
	for _, slug := range slugs {
		page, err := fetchPage(wikiSlug, slug)
		if err != nil {
			colors.RedPrint("Error fetching %s: %v", slug, err)
			continue
		}
		if err := savePage(*page); err != nil {
			colors.RedPrint("Failed to save page %s: %v", slug, err)
			continue
		}
		colors.GreenPrint("✅ Saved %s/%s.askr\n", wikiSlug, slug)
	}
}

func CloneWikis() {
	accessToken, err := auth.GetToken() // internal/auth の実装に合わせて GetAccessToken 等に適宜変更
	if err != nil {
		colors.RedPrint("Auth Error: %v", err)
		os.Exit(1)
	}

	resp, err := callAPIWikis(accessToken)
	if err != nil {
		colors.RedPrint("API Error: %s", err)
		os.Exit(1)
	}

	if len(resp) == 0 {
		colors.RedPrintText("No wikis found in Blob storage.")
		return
	}

	for _, wikiSlug := range resp {
		colors.GreenPrint("Cloning Wiki: %s...", wikiSlug)
		CloneWiki(wikiSlug)
	}
}
