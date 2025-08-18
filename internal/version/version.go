package version

import (
	"github.com/sakitibi/askrEditor/internal/colors"
)

const version = "2.0.20" // ビルド時に -ldflags で上書き可能

func PrintVersion() {
	colors.GreenPrint1("askreditor version", version)
}
