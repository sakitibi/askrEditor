package cmd

import (
	"os"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
	"github.com/sakitibi/askrEditor/internal/version"
	"github.com/sakitibi/askrEditor/internal/wiki"
)

func Execute() {
	if len(os.Args) < 2 {
		colors.RedPrintText("Usage: askreditor <clone|cloneall|push|login|version> ...")
		os.Exit(1)
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "clone":
		if len(os.Args) != 3 {
			colors.RedPrintText("Usage: askreditor clone <wikiSlug>")
			os.Exit(1)
			return
		}
		wiki.CloneWiki(os.Args[2])
	case "cloneall":
		if len(os.Args) != 2 {
			colors.RedPrintText("Usage askreditor cloneall")
			os.Exit(1)
		}
		wiki.CloneWikis()
	case "push":
		if len(os.Args) != 3 {
			colors.RedPrintText("Usage: askreditor push <wikiSlug>")
			os.Exit(1)
			return
		}
		PushWiki(os.Args[2])
	case "login":
		if len(os.Args) != 4 && len(os.Args) != 5 {
			colors.RedPrintText("Usage: askreditor login <email> <password>")
			os.Exit(1)
			return
		}
		var subCommand string
		if len(os.Args) > 4 {
			subCommand = os.Args[4]
		} else {
			subCommand = "default_value" // デフォルト値
		}
		if err := auth.Login(os.Args[2], os.Args[3], subCommand); err != nil {
			colors.RedPrint("Login error: %s", err)
			os.Exit(1)
		}
	case "version":
		version.PrintVersion()

	default:
		colors.RedPrint("Unknown command: %s", cmd)
		os.Exit(1)
	}
}
