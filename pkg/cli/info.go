package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	
	"github.com/allieus/imagekit/pkg/transform"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [image]",
	Short: "이미지 정보 표시",
	Long:  `이미지의 크기, 형식, DPI 등의 정보를 표시합니다.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func runInfo(cmd *cobra.Command, args []string) error {
	imagePath := args[0]
	
	// Open image file
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("파일을 열 수 없습니다: %w", err)
	}
	defer func() { _ = file.Close() }()
	
	// Read file content into buffer to allow multiple reads
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, file); err != nil {
		return fmt.Errorf("파일 읽기 실패: %w", err)
	}
	
	// Load image to get basic info
	img, format, err := transform.LoadImage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf("이미지 로드 실패: %w", err)
	}
	
	// Get image info
	info := transform.GetImageInfo(img, format)
	
	// Try to get DPI information
	dpi, _ := transform.GetImageDPI(bytes.NewReader(buf.Bytes()), format)
	if dpi > 0 {
		info.DPI = dpi
	}
	
	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}
	
	// Display information
	fmt.Println("📊 이미지 정보")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📁 파일명: %s\n", imagePath)
	fmt.Printf("📏 크기: %d x %d 픽셀\n", info.Width, info.Height)
	fmt.Printf("🎨 형식: %s\n", strings.ToUpper(string(info.Format)))
	fmt.Printf("📐 DPI: %d\n", info.DPI)
	fmt.Printf("💾 파일 크기: %s\n", formatFileSize(fileInfo.Size()))
	fmt.Printf("📅 수정 시간: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	return nil
}

func formatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}