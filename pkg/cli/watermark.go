package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	
	"github.com/allieus/pyhub-imagekit/pkg/transform"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	area   string
	method string
)

var watermarkCmd = &cobra.Command{
	Use:   "watermark [input] [output]",
	Short: "워터마크 제거",
	Long: `이미지에서 워터마크를 제거합니다.
	
예제:
  imagekit watermark --area=100,100,200,50 input.jpg output.jpg
  imagekit watermark --area=100,100,200,50 --method=blur input.jpg output.jpg`,
	Args: cobra.ExactArgs(2),
	RunE: runWatermark,
}

func init() {
	watermarkCmd.Flags().StringVar(&area, "area", "", "워터마크 영역 (x,y,width,height)")
	watermarkCmd.Flags().StringVar(&method, "method", "blur", "제거 방법 (blur, fill, inpaint)")
	watermarkCmd.MarkFlagRequired("area")
}

func runWatermark(cmd *cobra.Command, args []string) error {
	inputPath := args[0]
	outputPath := args[1]
	
	// Parse area
	rect, err := parseRectangle(area)
	if err != nil {
		return fmt.Errorf("영역 파싱 실패: %w", err)
	}
	
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
	bar := progressbar.Default(-1, "워터마크 제거 중...")
	
	// Remove watermark
	if err := transformer.RemoveWatermark(inputFile, outputFile, rect); err != nil {
		return fmt.Errorf("워터마크 제거 실패: %w", err)
	}
	
	bar.Finish()
	fmt.Printf("✅ 워터마크 제거 완료: %s\n", outputPath)
	
	return nil
}

func parseRectangle(area string) (transform.Rectangle, error) {
	parts := strings.Split(area, ",")
	if len(parts) != 4 {
		return transform.Rectangle{}, fmt.Errorf("영역은 x,y,width,height 형식이어야 합니다")
	}
	
	values := make([]int, 4)
	for i, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return transform.Rectangle{}, fmt.Errorf("잘못된 숫자: %s", part)
		}
		values[i] = val
	}
	
	return transform.Rectangle{
		X:      values[0],
		Y:      values[1],
		Width:  values[2],
		Height: values[3],
	}, nil
}