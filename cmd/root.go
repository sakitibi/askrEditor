package cmd

import (
	"os"

	"github.com/fatih/color"

	"github.com/sakitibi/askrEditor/internal/auth"
	"github.com/sakitibi/askrEditor/internal/colors"
	"github.com/sakitibi/askrEditor/internal/version"
	"github.com/sakitibi/askrEditor/internal/wiki"
)

func Execute() {
	if len(os.Args) < 2 {
		colors.RedPrintText("Usage: askreditor <clone|push|login|version> ...")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "clone":
		if len(os.Args) != 3 {
			colors.RedPrintText("Usage: askreditor clone <wikiSlug>")
			return
		}
		wiki.CloneWiki(os.Args[2])
	case "cloneall":
		if len(os.Args) != 2 {
			colors.RedPrintText("Usage askreditor cloneall")
		}
		wiki.CloneWikis()
	case "push":
		if len(os.Args) != 3 {
			colors.RedPrintText("Usage: askreditor push <wikiSlug>")
			return
		}
		PushWiki(os.Args[2])
	case "login":
		if len(os.Args) != 4 {
			colors.RedPrintText("Usage: askreditor login <email> <password>")
			return
		}
		if err := auth.Login(os.Args[2], os.Args[3]); err != nil {
			color.New(color.FgRed, color.Bold).Println("Login error:", err)
		}
	case "version":
		version.PrintVersion()

	default:
		colors.RedPrint("Unknown command: %s", cmd)
	}
}
