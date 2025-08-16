package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	supabase "github.com/supabase-community/supabase-go"
)

func ReplaceFile(client *supabase.Client, file string) error {
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
		Update(updates, "public", "wiki_pages").
		Eq("slug", slug).
		Execute()
	if err != nil {
		return err
	}

	fmt.Printf("✅ Replaced in public.wiki_pages where slug='%s'\n", slug)
	return nil
}
