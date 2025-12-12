package version

import (
	"github.com/sakitibi/askrEditor/internal/colors"
)

const version = "2.0.60" // ビルド時に -ldflags で上書き可能

func PrintVersion() {
	colors.GreenPrint("askreditor version%s", version)
}
