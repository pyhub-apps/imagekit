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
	width   string
	height  string
	dpi     int
	mode    string
	quality int
)

var convertCmd = &cobra.Command{
	Use:   "convert [input-pattern or file] [output-file (optional)]",
	Short: "이미지 변환 (크기, DPI)",
	Long: `단일 파일 또는 glob 패턴으로 여러 이미지를 변환합니다.
	
예제:
  # 단일 파일 변환
  imagekit convert --width=1920 --height=1080 input.jpg output.jpg
  imagekit convert --dpi=96 input.png output.png
  
  # 여러 파일 변환 (glob 패턴)
  imagekit convert --width=1920 "*.jpg"              # 모든 jpg 파일
  imagekit convert --dpi=96 "photos/*.png"           # photos 디렉토리의 png 파일들
  imagekit convert --width=800 --height=600 "*.{jpg,png}"  # jpg와 png 파일들`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runConvert,
}

func init() {
	convertCmd.Flags().StringVar(&width, "width", "", "목표 너비 (픽셀 또는 배수: 1920, 2x, x2, 0.5x)")
	convertCmd.Flags().StringVar(&height, "height", "", "목표 높이 (픽셀 또는 배수: 1080, 2x, x2, 0.5x)")
	convertCmd.Flags().IntVar(&dpi, "dpi", 0, "목표 DPI (72, 96, 150, 300)")
	convertCmd.Flags().StringVar(&mode, "mode", "fit", "리사이징 모드 (fit, fill, exact)")
	convertCmd.Flags().IntVar(&quality, "quality", 95, "JPEG 품질 (1-100)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputPattern := args[0]
	
	// Check if conversion options are specified
	if width == "" && height == "" && dpi <= 0 {
		return fmt.Errorf("변환 옵션을 지정해주세요 (--width, --height, 또는 --dpi)")
	}
	
	// Create transformer
	transformer := transform.NewTransformer()
	
	// Check if it's a glob pattern or contains wildcards
	hasGlob := strings.Contains(inputPattern, "*") || strings.Contains(inputPattern, "?") || strings.Contains(inputPattern, "[")
	
	// Single file mode with explicit output
	if len(args) == 2 && !hasGlob {
		outputPath := args[1]
		return processSingleFile(transformer, inputPattern, outputPath)
	}
	
	// Check if it's a single file without glob patterns
	if !hasGlob {
		// Single file mode with auto-generated output name
		if _, err := os.Stat(inputPattern); err == nil {
			outputPath := batch.GenerateOutputPath(inputPattern)
			return processSingleFile(transformer, inputPattern, outputPath)
		}
		return fmt.Errorf("파일을 찾을 수 없습니다: %s", inputPattern)
	}
	
	// Batch mode
	processor := batch.NewProcessor(transformer)
	
	// Parse dimensions
	widthDim, err := transform.ParseDimension(width)
	if err != nil {
		return fmt.Errorf("잘못된 width 값: %w", err)
	}
	heightDim, err := transform.ParseDimension(height)
	if err != nil {
		return fmt.Errorf("잘못된 height 값: %w", err)
	}
	
	// Prepare options
	var resizeOptions *transform.ResizeOptions
	if !widthDim.IsZero() || !heightDim.IsZero() {
		resizeMode := getResizeMode(mode)
		resizeOptions = &transform.ResizeOptions{
			WidthDim:  widthDim,
			HeightDim: heightDim,
			Mode:      resizeMode,
			Quality:   quality,
		}
	}
	
	options := batch.ProcessOptions{
		ResizeOptions: resizeOptions,
		DPI:           dpi,
	}
	
	// Progress callback
	fmt.Println("Converting images...")
	progressCallback := func(current, total int, fileName string, success bool) {
		status := "✅"
		if !success {
			status = "❌"
		}
		fmt.Printf("[%d/%d] %s → %s %s\n", current, total, fileName, 
			strings.TrimSuffix(fileName, filepath.Ext(fileName))+"_converted"+filepath.Ext(fileName), status)
	}
	
	// Process files
	result, err := processor.ProcessFiles(inputPattern, options, progressCallback)
	if err != nil {
		return err
	}
	
	// Show summary
	fmt.Printf("\n완료: %d/%d 성공", result.SuccessCount, result.TotalFiles)
	if result.HasErrors() {
		fmt.Printf(", %d 실패\n", len(result.FailedFiles))
		fmt.Println("\n실패한 파일:")
		for _, failed := range result.FailedFiles {
			fmt.Printf("  - %s: %v\n", failed.Path, failed.Error)
		}
	} else {
		fmt.Println()
	}
	
	return nil
}

// processSingleFile handles single file conversion
func processSingleFile(transformer *transform.Transformer, inputPath, outputPath string) error {
	// Show progress
	bar := progressbar.Default(-1, "이미지 변환 중...")
	
	// Parse dimensions
	widthDim, err := transform.ParseDimension(width)
	if err != nil {
		return fmt.Errorf("잘못된 width 값: %w", err)
	}
	heightDim, err := transform.ParseDimension(height)
	if err != nil {
		return fmt.Errorf("잘못된 height 값: %w", err)
	}
	
	// Process based on flags
	if !widthDim.IsZero() || !heightDim.IsZero() {
		// Open input file for resize
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("입력 파일을 열 수 없습니다: %w", err)
		}
		defer func() { _ = inputFile.Close() }()
		
		// Create output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("출력 파일을 생성할 수 없습니다: %w", err)
		}
		defer func() { _ = outputFile.Close() }()
		
		// Resize operation
		resizeMode := getResizeMode(mode)
		options := transform.ResizeOptions{
			WidthDim:  widthDim,
			HeightDim: heightDim,
			Mode:      resizeMode,
			Quality:   quality,
		}
		
		bar.Describe("크기 변환 중...")
		if err := transformer.Resize(inputFile, outputFile, options); err != nil {
			return fmt.Errorf("크기 변환 실패: %w", err)
		}
		
		// If DPI is also specified, we need to process it separately
		if dpi > 0 {
			// Re-open files for DPI processing
			_ = inputFile.Close()
			_ = outputFile.Close()
			
			// Use the resized output as input for DPI change
			tempFile, err := os.Open(outputPath)
			if err != nil {
				return fmt.Errorf("임시 파일을 열 수 없습니다: %w", err)
			}
			defer func() { _ = tempFile.Close() }()
			
			outputFile2, err := os.Create(outputPath + ".tmp")
			if err != nil {
				return fmt.Errorf("임시 출력 파일을 생성할 수 없습니다: %w", err)
			}
			defer func() { _ = outputFile2.Close() }()
			
			bar.Describe("DPI 설정 중...")
			if err := transformer.SetDPI(tempFile, outputFile2, dpi); err != nil {
				return fmt.Errorf("DPI 설정 실패: %w", err)
			}
			
			// Replace original with DPI-adjusted version
			if err := os.Rename(outputPath+".tmp", outputPath); err != nil {
				return fmt.Errorf("failed to rename temporary file: %w", err)
			}
		}
	} else if dpi > 0 {
		// Open input file for DPI
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("입력 파일을 열 수 없습니다: %w", err)
		}
		defer func() { _ = inputFile.Close() }()
		
		// Create output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("출력 파일을 생성할 수 없습니다: %w", err)
		}
		defer func() { _ = outputFile.Close() }()
		
		// DPI only operation
		bar.Describe("DPI 설정 중...")
		if err := transformer.SetDPI(inputFile, outputFile, dpi); err != nil {
			return fmt.Errorf("DPI 설정 실패: %w", err)
		}
	}
	
	_ = bar.Finish()
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