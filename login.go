package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

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

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func login(email, password string) error {
	// Supabase の API エンドポイント
	url := supabaseURL + "/auth/v1/token?grant_type=password"

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
	req.Header.Set("apikey", supabaseAnonKey)
	req.Header.Set("Content-Type", "application/json")

	// 実行
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", body)
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

	fmt.Println("✅ Login successful, tokens saved")
	return nil
}

func getToken() (string, error) {
	home, _ := os.UserHomeDir()
	tokenPath := filepath.Join(home, ".askreditor_token.json")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("not logged in, please run `askreditor login <email> <password>`")
	}
	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return "", fmt.Errorf("invalid token file, try logging in again")
	}
	return token.AccessToken, nil
}
