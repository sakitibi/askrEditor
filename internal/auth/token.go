package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".askreditor_token.json")
}

func SaveToken(result LoginResponse) error {
	data, _ := json.MarshalIndent(result, "", "  ")
	return os.WriteFile(tokenPath(), data, 0600)
}

func GetToken() (string, error) {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return "", fmt.Errorf("not logged in, please run `askreditor login <email> <password>`")
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(data, &token); err != nil {
		return "", fmt.Errorf("invalid token file, try logging in again")
	}
	return token.AccessToken, nil
}
