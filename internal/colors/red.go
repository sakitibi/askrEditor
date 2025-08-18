package colors

import "github.com/fatih/color"

func RedPrint(content string, args1 any, args2 any) {
	color.New(color.FgRed, color.Bold).Printf(content, args1, args2)
}

func RedPrint1(content string, args1 any) {
	color.New(color.FgRed, color.Bold).Printf(content, args1)
}

func RedPrintText(content string) {
	color.New(color.FgRed, color.Bold).Printf(content)
}
