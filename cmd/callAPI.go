package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sakitibi/askrEditor/internal/auth"
)

func callAPI(method, wikiSlug, pageSlug string, body map[string]string, accessToken string, apiKey string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%s", auth.ApiBaseURL, wikiSlug, pageSlug)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// ヘッダー設定
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	// apiKey が渡された場合のみ x-apikey ヘッダーを設定
	if apiKey != "" {
		req.Header.Set("x-apikey", apiKey)
	}

	// サーバー側で CLI からのアクセスを判別するためのカスタムヘッダー
	req.Header.Set("X-CLI", "true")

	// タイムアウトを設定したクライアントを使用することを推奨します
	client := &http.Client{
		Timeout: 40 * time.Second,
	}

	return client.Do(req)
}
