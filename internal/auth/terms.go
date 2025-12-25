package auth

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sakitibi/askrEditor/internal/colors"
)

func CheckTerms() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	termsFile := filepath.Join(home, ".terms_agreed")

	// すでに存在する場合は何もしない
	if _, err := os.Stat(termsFile); err == nil {
		return nil
	}

	// デフォルトブラウザで URL を開く
	url := "https://asakura-wiki.vercel.app/policies" // 実際の利用規約URLに変更
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // Linux
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		colors.RedPrint("Failed to open browser:%s", err)
	}

	// プロンプト表示
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you agree to the terms? (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		// ファイル作成して同意を記録
		if err := os.WriteFile(termsFile, []byte("agreed"), 0600); err != nil {
			return err
		}
		colors.GreenPrintText("Terms agreed")
		return nil
	} else {
		colors.RedPrintText("You must agree to the terms to continue")
		os.Exit(1)
	}
	return nil
}
