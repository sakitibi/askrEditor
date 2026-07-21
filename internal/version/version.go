package version

import (
	"github.com/sakitibi/askrEditor/internal/colors"
)

const version = "2.1.0" // ビルド時に -ldflags で上書き可能

func PrintVersion() {
	colors.GreenPrint("askreditor version%s", version)
	colors.GreenPrintText("CopyRight 2025~2026 14ninstudio, Inc All Rights Reserved.")
}
