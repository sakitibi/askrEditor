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
		colors.GreenPrintText("以下に認証用の問題:  ")
		key, value, _ := GetRandomEntry(MapactionsFunc())
		colors.GreenPrintText(key)
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			input := scanner.Text()
			if input == value {
				colors.GreenPrintText("正解!")
			} else {
				colors.RedPrintText("不正解..")
				os.Exit(1)
			}

			// 分離した関数を呼び出す
			if err := loginAndSaveToken(email, password); err != nil {
				os.Exit(1)
			}
			return nil
		}
	}
	return nil
}

// 分離した関数: Supabaseでの認証とトークンの保存を行う
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
