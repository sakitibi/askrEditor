package main

import (
	"fmt"
	"os"
)

func printError(msg string, err any) {
	// 赤文字で出力
	fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", msg)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: askreditor <command> [args...]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "clone":
		if len(os.Args) < 3 {
			printError("Usage: askreditor clone <wiki_slug>", "")
			os.Exit(1)
		}
		wikiSlug := os.Args[2]

		client, err := NewSupabaseClient()
		if err != nil {
			fmt.Println("Cannot initialize Supabase client:", err)
			os.Exit(1)
		}

		if err := cloneWiki(client, wikiSlug); err != nil {
			printError("Clone failed:", err)
			os.Exit(1)
		}

	case "replace":
		if len(os.Args) < 3 {
			fmt.Println("Usage: askreditor replace <file.askr>")
			os.Exit(1)
		}
		file := os.Args[2]

		client, err := NewSupabaseClient()
		if err != nil {
			fmt.Println("Cannot initialize Supabase client:", err)
			os.Exit(1)
		}

		if err := ReplaceFile(client, file); err != nil {
			fmt.Println("Replace failed:", err)
			os.Exit(1)
		}

	case "version":
		PrintVersion()

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
