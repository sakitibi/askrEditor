package colors

import "github.com/fatih/color"

func RedPrint(format string, a ...any) {
	color.New(color.FgRed, color.Bold).Printf(format+"\n", a...)
}

func RedPrintText(content string) {
	color.New(color.FgRed, color.Bold).Println(content)
}
