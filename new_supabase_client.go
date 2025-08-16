package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	supabase "github.com/supabase-community/supabase-go"
)

func init() {
	// .env を読み込む（なければ無視）
	_ = godotenv.Load()
}

func NewSupabaseClient() (*supabase.Client, error) {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseUrl == "" || supabaseKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set")
	}

	client, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
