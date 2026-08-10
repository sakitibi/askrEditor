package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sakitibi/askrEditor/internal/colors"
)

const ApiBaseURL = "https://asakura-wiki.vercel.app/api/wiki_v2"

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func Login(email, password string, skipFlug string) error {
	// 利用規約確認
	var skipFlugBool bool = false
	if skipFlug == "--skip" {
		skipFlugBool = true
	}
	if err := CheckTerms(); err != nil {
		return err
	}
	if skipFlugBool {
		if err := loginAndSaveToken(email, password); err != nil {
			os.Exit(1)
		}
		return nil
	} else {
		actionMap := MapactionsFunc()

		// 1問目の取得
		key1, value1, ok1 := GetRandomEntry(actionMap)
		if !ok1 {
			return fmt.Errorf("認証問題の取得に失敗しました")
		}

		remainingMap := make(map[string]string)
		for k, v := range actionMap {
			if k != key1 {
				remainingMap[k] = v
			}
		}

		// 2問目の取得
		key2, value2, ok2 := GetRandomEntry(remainingMap)
		if !ok2 {
			return fmt.Errorf("2問目の認証問題の取得に失敗しました")
		}

		questions := []struct {
			key   string
			value string
		}{
			{key: key1, value: value1},
			{key: key2, value: value2},
		}

		// スキャナーをループの外で1回だけ生成
		scanner := bufio.NewScanner(os.Stdin)

		// 2問連続で出題
		for i, q := range questions {
			colors.GreenPrintText(fmt.Sprintf("以下に認証用の問題 (%d/2):", i+1))
			colors.GreenPrintText(q.key)

			if scanner.Scan() {
				input := scanner.Text()
				if input == q.value {
					colors.GreenPrintText("正解!")
				} else {
					colors.RedPrintText("不正解..")
					os.Exit(1)
				}
			} else {
				if err := scanner.Err(); err != nil {
					colors.RedPrint("入力エラーが発生しました: %v", err)
				}
				os.Exit(1)
			}
		}

		// 全問正解後にログイン実行
		if err := loginAndSaveToken(email, password); err != nil {
			os.Exit(1)
		}
		return nil
	}
}

func loginAndSaveToken(email, password string) error {
	// Supabase の API エンドポイント
	url := SupabaseURL + "/auth/v1/token?grant_type=password"

	// リクエストボディ
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	// JSON に変換
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// POST リクエスト作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("apikey", SupabaseAnonKey)
	req.Header.Set("Content-Type", "application/json")

	// 実行
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		colors.RedPrint("login failed: %s", body)
		return fmt.Errorf("login failed")
	}

	// レスポンスを構造体にパース
	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// 保存先
	configDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, ".askreditor_token.json")

	// 保存する
	tokenData, _ := json.MarshalIndent(result, "", "  ")
	if err := os.WriteFile(configFile, tokenData, 0600); err != nil {
		return err
	}

	colors.GreenPrintText("✅ Login successful, tokens saved")
	return nil
}
