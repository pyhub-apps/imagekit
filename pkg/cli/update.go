package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/allieus/pyhub-imagekit/pkg/update"
	"github.com/spf13/cobra"
)

var (
	checkOnly bool
	forceUpdate bool
	targetVersion string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "imagekit 최신 버전으로 업데이트",
	Long: `imagekit을 최신 버전으로 업데이트합니다.
	
예제:
  imagekit update              # 최신 버전으로 업데이트
  imagekit update --check       # 업데이트 가능 여부만 확인
  imagekit update --force       # 강제 재설치`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "업데이트 가능 여부만 확인")
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false, "현재 버전과 동일해도 강제 재설치")
	updateCmd.Flags().StringVar(&targetVersion, "version", "", "특정 버전으로 업데이트")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := rootCmd.Version
	
	updater, err := update.NewUpdater(currentVersion)
	if err != nil {
		return fmt.Errorf("업데이터 초기화 실패: %w", err)
	}
	
	fmt.Printf("현재 버전: %s\n", currentVersion)
	fmt.Println("최신 버전 확인 중...")
	
	release, hasUpdate, err := updater.CheckForUpdate()
	if err != nil {
		return fmt.Errorf("버전 확인 실패: %w", err)
	}
	
	fmt.Printf("최신 버전: %s\n", release.TagName)
	
	if checkOnly {
		if hasUpdate {
			fmt.Println("\n✨ 새 버전이 있습니다!")
			fmt.Printf("변경사항: %s\n", release.HTMLURL)
		} else {
			fmt.Println("\n✅ 최신 버전을 사용 중입니다.")
		}
		return nil
	}
	
	if !hasUpdate && !forceUpdate {
		fmt.Println("\n✅ 이미 최신 버전을 사용 중입니다.")
		return nil
	}
	
	// Show changes
	if release.Body != "" {
		fmt.Println("\n변경사항:")
		summary := release.GetChangesSummary()
		fmt.Println(summary)
	}
	
	// Confirm update
	if !confirmUpdate() {
		fmt.Println("업데이트가 취소되었습니다.")
		return nil
	}
	
	// Perform update
	fmt.Println("\n업데이트를 시작합니다...")
	if err := updater.Update(release, forceUpdate); err != nil {
		return fmt.Errorf("업데이트 실패: %w", err)
	}
	
	fmt.Println("✅ 업데이트가 완료되었습니다!")
	fmt.Println("새 버전을 확인하려면 다음 명령을 실행하세요: imagekit --version")
	
	return nil
}

func confirmUpdate() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n업데이트를 진행하시겠습니까? (y/N): ")
	
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}