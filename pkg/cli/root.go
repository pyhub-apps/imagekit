package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "imagekit",
	Short: "이미지 변환 CLI 도구",
	Long: `imagekit은 이미지 최적화 도구입니다.
	
JPG 및 PNG 이미지의 크기, DPI를 변환하고 워터마크를 제거할 수 있습니다.`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(watermarkCmd)
	rootCmd.AddCommand(infoCmd)
}
