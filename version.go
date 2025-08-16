package main

import "fmt"

const version = "2.0.20" // ビルド時に -ldflags で上書き可能

func PrintVersion() {
	fmt.Println("askreditor version", version)
}
