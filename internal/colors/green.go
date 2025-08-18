package colors

import "github.com/fatih/color"

func GreenPrint(content string, args1 any, args2 any) {
	color.New(color.FgGreen, color.Bold).Println(content, args1, args2)
}

func GreenPrint1(content string, args1 any) {
	color.New(color.FgGreen, color.Bold).Println(content, args1)
}

func GreenPrintText(content string) {
	color.New(color.FgGreen, color.Bold).Println(content)
}
