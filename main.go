package main

import (
	"fmt"
	"os"
)

const apiBaseURL = "https://asakura-wiki.vercel.app/api/wiki"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: askreditor <replace|clone|delete|get|push> ...")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "replace":
		if len(os.Args) < 5 {
			fmt.Println("Usage: askreditor replace <wikiSlug> <pageSlug> <content>")
			return
		}
		wikiSlug := os.Args[2]
		pageSlug := os.Args[3]
		content := os.Args[4]
		callAPI("PUT", wikiSlug, pageSlug, map[string]string{"content": content}, "true") // X-CLI=true

	case "clone":
		if len(os.Args) < 3 {
			fmt.Println("Usage: askreditor clone <wikiSlug>")
			return
		}
		wikiSlug := os.Args[2]
		cloneWiki(wikiSlug)

	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Usage: askreditor delete <wikiSlug> <pageSlug>")
			return
		}
		wikiSlug := os.Args[2]
		pageSlug := os.Args[3]
		callAPI("DELETE", wikiSlug, pageSlug, nil, "true")

	case "get":
		if len(os.Args) < 4 {
			fmt.Println("Usage: askreditor get <wikiSlug> <pageSlug>")
			return
		}
		wikiSlug := os.Args[2]
		pageSlug := os.Args[3]
		callAPI("GET", wikiSlug, pageSlug, nil, "")

	case "push":
		if len(os.Args) < 3 {
			fmt.Println("Usage: askreditor push <wikiSlug>")
			return
		}
		wikiSlug := os.Args[2]
		pushWiki(wikiSlug)

	case "version":
		PrintVersion()
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
