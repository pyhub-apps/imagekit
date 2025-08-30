package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/allieus/pyhub-imagekit/pkg/batch"
	"github.com/allieus/pyhub-imagekit/pkg/transform"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	cropTop    string
	cropBottom string
	cropLeft   string
	cropRight  string
)

var cropCmd = &cobra.Command{
	Use:   "crop [input-pattern or file] [output-file (optional)]",
	Short: "이미지 가장자리 크롭",
	Long: `이미지의 가장자리를 잘라냅니다. 워터마크 제거나 여백 제거에 유용합니다.
	
예제:
  # 단일 파일 크롭
  imagekit crop --bottom=100 watermarked.jpg clean.jpg        # 하단 100픽셀 제거
  imagekit crop --top=10% --bottom=5% input.jpg output.jpg    # 상단 10%, 하단 5% 제거
  
  # 여러 파일 크롭 (glob 패턴)
  imagekit crop --bottom=50 "*.jpg"                           # 모든 jpg 파일 하단 50픽셀 제거
  imagekit crop --top=15% "photos/*.png"                      # photos 디렉토리의 png 파일들 상단 15% 제거
  
  # 모든 가장자리 크롭
  imagekit crop --top=20 --bottom=20 --left=20 --right=20 input.jpg output.jpg`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runCrop,
}

func init() {
	cropCmd.Flags().StringVar(&cropTop, "top", "", "상단에서 제거할 영역 (픽셀 또는 %)")
	cropCmd.Flags().StringVar(&cropBottom, "bottom", "", "하단에서 제거할 영역 (픽셀 또는 %)")
	cropCmd.Flags().StringVar(&cropLeft, "left", "", "좌측에서 제거할 영역 (픽셀 또는 %)")
	cropCmd.Flags().StringVar(&cropRight, "right", "", "우측에서 제거할 영역 (픽셀 또는 %)")
}

func runCrop(cmd *cobra.Command, args []string) error {
	inputPattern := args[0]
	
	// Check if at least one crop option is specified
	if cropTop == "" && cropBottom == "" && cropLeft == "" && cropRight == "" {
		return fmt.Errorf("최소 하나의 크롭 옵션을 지정해주세요 (--top, --bottom, --left, --right)")
	}
	
	// Parse crop options
	options, err := parseCropOptions()
	if err != nil {
		return fmt.Errorf("크롭 옵션 파싱 실패: %w", err)
	}
	
	// Create transformer
	transformer := transform.NewTransformer()
	
	// Check if it's a glob pattern or contains wildcards
	hasGlob := strings.Contains(inputPattern, "*") || strings.Contains(inputPattern, "?") || strings.Contains(inputPattern, "[")
	
	// Single file mode with explicit output
	if len(args) == 2 && !hasGlob {
		outputPath := args[1]
		return processSingleCropFile(transformer, inputPattern, outputPath, options)
	}
	
	// Check if it's a single file without glob patterns
	if !hasGlob {
		// Single file mode with auto-generated output name
		if _, err := os.Stat(inputPattern); err == nil {
			outputPath := batch.GenerateOutputPath(inputPattern)
			return processSingleCropFile(transformer, inputPattern, outputPath, options)
		}
		return fmt.Errorf("파일을 찾을 수 없습니다: %s", inputPattern)
	}
	
	// Batch mode
	return processBatchCrop(transformer, inputPattern, options)
}

func parseCropOptions() (transform.EdgeCropOptions, error) {
	options := transform.EdgeCropOptions{}
	
	if cropTop != "" {
		top, err := transform.ParseCropValue(cropTop)
		if err != nil {
			return options, fmt.Errorf("잘못된 top 값: %w", err)
		}
		options.Top = top
	}
	
	if cropBottom != "" {
		bottom, err := transform.ParseCropValue(cropBottom)
		if err != nil {
			return options, fmt.Errorf("잘못된 bottom 값: %w", err)
		}
		options.Bottom = bottom
	}
	
	if cropLeft != "" {
		left, err := transform.ParseCropValue(cropLeft)
		if err != nil {
			return options, fmt.Errorf("잘못된 left 값: %w", err)
		}
		options.Left = left
	}
	
	if cropRight != "" {
		right, err := transform.ParseCropValue(cropRight)
		if err != nil {
			return options, fmt.Errorf("잘못된 right 값: %w", err)
		}
		options.Right = right
	}
	
	return options, nil
}

func processSingleCropFile(transformer *transform.Transformer, inputPath, outputPath string, options transform.EdgeCropOptions) error {
	// Show progress
	bar := progressbar.Default(-1, "이미지 크롭 중...")
	
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
	
	// Perform crop
	if err := transformer.CropEdges(inputFile, outputFile, options); err != nil {
		return fmt.Errorf("크롭 실패: %w", err)
	}
	
	bar.Finish()
	fmt.Printf("✅ 크롭 완료: %s\n", outputPath)
	
	return nil
}

func processBatchCrop(transformer *transform.Transformer, pattern string, options transform.EdgeCropOptions) error {
	// Find matching files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("잘못된 glob 패턴: %w", err)
	}
	
	if len(matches) == 0 {
		return fmt.Errorf("패턴과 일치하는 파일이 없습니다: %s", pattern)
	}
	
	// Filter valid image files
	var filesToProcess []string
	for _, match := range matches {
		// Skip already converted files
		if batch.IsConvertedFile(match) {
			continue
		}
		
		// Check if it's a supported image format
		ext := strings.ToLower(filepath.Ext(match))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}
		
		filesToProcess = append(filesToProcess, match)
	}
	
	if len(filesToProcess) == 0 {
		return fmt.Errorf("처리할 유효한 이미지 파일이 없습니다")
	}
	
	// Process files
	fmt.Println("Cropping images...")
	successCount := 0
	var failedFiles []string
	
	for i, inputPath := range filesToProcess {
		outputPath := batch.GenerateOutputPath(inputPath)
		
		// Process single file
		err := processSingleCropFile(transformer, inputPath, outputPath, options)
		
		status := "✅"
		if err != nil {
			status = "❌"
			failedFiles = append(failedFiles, inputPath)
		} else {
			successCount++
		}
		
		fmt.Printf("[%d/%d] %s → %s %s\n", i+1, len(filesToProcess), 
			filepath.Base(inputPath), filepath.Base(outputPath), status)
		
		if err != nil {
			fmt.Printf("  에러: %v\n", err)
		}
	}
	
	// Show summary
	fmt.Printf("\n완료: %d/%d 성공", successCount, len(filesToProcess))
	if len(failedFiles) > 0 {
		fmt.Printf(", %d 실패\n", len(failedFiles))
		fmt.Println("\n실패한 파일:")
		for _, path := range failedFiles {
			fmt.Printf("  - %s\n", path)
		}
	} else {
		fmt.Println()
	}
	
	return nil
}