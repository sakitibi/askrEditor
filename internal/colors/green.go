package colors

import "github.com/fatih/color"

func GreenPrint(format string, a ...any) {
	color.New(color.FgGreen, color.Bold).Printf(format+"\n", a...)
}

func GreenPrintText(content string) {
	color.New(color.FgGreen, color.Bold).Println(content)
}
