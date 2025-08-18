package wiki

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
)

type Page struct {
	Slug     string `json:"slug"`
	WikiSlug string `json:"wiki_slug"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// savePage は page を wikiSlug/slug.askr に保存
func savePage(page Page) error {
	filePath := filepath.Join(page.WikiSlug, page.Slug+".askr")
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// TITLE 行を追加
	data := fmt.Sprintf("TITLE:%s\n%s", page.Title, page.Content)
	return os.WriteFile(filePath, []byte(data), 0644)
}

// fetchPage は API からページを取得
func fetchPage(wikiSlug, pageSlug string) (*Page, error) {
	url := fmt.Sprintf("https://asakura-wiki.vercel.app/api/wiki/%s/%s", wikiSlug, pageSlug)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var page Page
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &page, nil
}

// fetchSlugs は API から wikiSlug のページ一覧を取得
func fetchSlugs(wikiSlug string) ([]string, error) {
	url := fmt.Sprintf("%s/%s", auth.ApiBaseURL, wikiSlug)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch slugs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// page_slugs を取り出す
	var data struct {
		WikiSlug  string   `json:"wiki_slug"`
		PageSlugs []string `json:"page_slugs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data.PageSlugs, nil
}

// CloneWiki は wikiSlug を指定して全ページをローカルに保存
func CloneWiki(wikiSlug string) {
	slugs, err := fetchSlugs(wikiSlug)
	if err != nil {
		colors.RedPrint("Failed to fetch slug list: %s", err)
		return
	}
	if len(slugs) == 0 {
		colors.RedPrint("%s is Not defined", wikiSlug)
		return
	}
	for _, slug := range slugs {
		page, err := fetchPage(wikiSlug, slug)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := savePage(*page); err != nil {
			colors.RedPrint("Failed to save page: %s", err)
			continue
		}
		colors.GreenPrint("✅ Saved %s/%s.askr\n", page.WikiSlug, page.Slug)
	}
}
