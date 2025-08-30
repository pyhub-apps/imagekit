package cli

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/allieus/pyhub-imagekit/pkg/transform"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	width   int
	height  int
	dpi     int
	mode    string
	quality int
)

var convertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "이미지 변환 (크기, DPI)",
	Long: `이미지의 크기와 DPI를 변환합니다.
	
예제:
  imagekit convert --width=1920 --height=1080 input.jpg output.jpg
  imagekit convert --dpi=96 input.png output.png
  imagekit convert --width=800 --mode=fit input.jpg output.jpg`,
	Args: cobra.ExactArgs(2),
	RunE: runConvert,
}

func init() {
	convertCmd.Flags().IntVar(&width, "width", 0, "목표 너비 (픽셀)")
	convertCmd.Flags().IntVar(&height, "height", 0, "목표 높이 (픽셀)")
	convertCmd.Flags().IntVar(&dpi, "dpi", 0, "목표 DPI (72, 96, 150, 300)")
	convertCmd.Flags().StringVar(&mode, "mode", "fit", "리사이징 모드 (fit, fill, exact)")
	convertCmd.Flags().IntVar(&quality, "quality", 95, "JPEG 품질 (1-100)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputPath := args[0]
	outputPath := args[1]
	
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("입력 파일을 열 수 없습니다: %w", err)
	}
	defer inputFile.Close()
	
	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("출력 파일을 생성할 수 없습니다: %w", err)
	}
	defer outputFile.Close()
	
	// Create transformer
	transformer := transform.NewTransformer()
	
	// Show progress
	bar := progressbar.Default(-1, "이미지 변환 중...")
	
	// Process based on flags
	if width > 0 || height > 0 {
		// Resize operation
		resizeMode := getResizeMode(mode)
		options := transform.ResizeOptions{
			Width:   width,
			Height:  height,
			Mode:    resizeMode,
			Quality: quality,
		}
		
		bar.Describe("크기 변환 중...")
		if err := transformer.Resize(inputFile, outputFile, options); err != nil {
			return fmt.Errorf("크기 변환 실패: %w", err)
		}
		
		// If DPI is also specified, we need to process it separately
		if dpi > 0 {
			// Re-open files for DPI processing
			inputFile.Close()
			outputFile.Close()
			
			// Use the resized output as input for DPI change
			tempFile, err := os.Open(outputPath)
			if err != nil {
				return fmt.Errorf("임시 파일을 열 수 없습니다: %w", err)
			}
			defer tempFile.Close()
			
			outputFile2, err := os.Create(outputPath + ".tmp")
			if err != nil {
				return fmt.Errorf("임시 출력 파일을 생성할 수 없습니다: %w", err)
			}
			defer outputFile2.Close()
			
			bar.Describe("DPI 설정 중...")
			if err := transformer.SetDPI(tempFile, outputFile2, dpi); err != nil {
				return fmt.Errorf("DPI 설정 실패: %w", err)
			}
			
			// Replace original with DPI-adjusted version
			os.Rename(outputPath+".tmp", outputPath)
		}
	} else if dpi > 0 {
		// DPI only operation
		bar.Describe("DPI 설정 중...")
		if err := transformer.SetDPI(inputFile, outputFile, dpi); err != nil {
			return fmt.Errorf("DPI 설정 실패: %w", err)
		}
	} else {
		return fmt.Errorf("변환 옵션을 지정해주세요 (--width, --height, 또는 --dpi)")
	}
	
	bar.Finish()
	fmt.Printf("✅ 변환 완료: %s\n", outputPath)
	
	return nil
}

func getResizeMode(mode string) transform.ResizeMode {
	switch strings.ToLower(mode) {
	case "fill":
		return transform.ResizeFill
	case "exact":
		return transform.ResizeExact
	default:
		return transform.ResizeFit
	}
}