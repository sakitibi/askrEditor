package main

import (
	"fmt"
	"io/ioutil"
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: askr replace <file.askr>")
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

	if command == "replace" {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		// ファイル名 (拡張子 .askr を除去)
		slug := file[:len(file)-5]

		_, _, err = client.
			From("wiki_pages").
			Update(
				map[string]interface{}{
					"content": string(data),
				},
				"public",
				"wiki_pages",
			).
			Eq("slug", slug).
			Execute()

		if err != nil {
			fmt.Println("Replace failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Replaced in public.wiki_pages where slug='%s'\n", slug)
	} else {
		fmt.Println("Unknown command:", command)
	}
}
