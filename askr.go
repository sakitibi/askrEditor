package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	supabase "github.com/supabase-community/supabase-go"
)

const version = "2.0.20" // ← ビルド時に書き換え可能

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: askr <command> [args...]")
		os.Exit(1)
	}
	command := os.Args[1]
	file := os.Args[2]

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseUrl == "" || supabaseKey == "" {
		fmt.Println("Error: SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set")
		os.Exit(1)
	}

	client, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)
	if err != nil {
		fmt.Println("Cannot initialize Supabase client:", err)
		os.Exit(1)
	}

	switch command {
	case "replace":
		err := replaceFile(client, file)
		if err != nil {
			fmt.Println("Replace failed:", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("askreditor version", version)
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}

func replaceFile(client *supabase.Client, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	slug := strings.TrimSuffix(filepath.Base(file), ".askr")

	updates := map[string]interface{}{
		"content": string(data),
	}

	_, _, err = client.
		From("wiki_pages").
		Update(updates, "public", "wiki_pages"). // <- Update は現在 interface{}, string, string の3引数
		Eq("slug", slug).
		Execute()
	if err != nil {
		return err
	}

	fmt.Printf("✅ Replaced in public.wiki_pages where slug='%s'\n", slug)
	return nil
}
