package colors

import "github.com/fatih/color"

func GreenPrint(format string, args ...any) {
	color.New(color.FgGreen, color.Bold).Printf(format, args...)
}

func GreenPrint1(content string, args1 any) {
	color.New(color.FgGreen, color.Bold).Printf(content, args1)
}

func GreenPrintText(content string) {
	color.New(color.FgGreen, color.Bold).Printf(content)
}
