package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Updater handles the self-update process
type Updater struct {
	currentVersion string
	client         *GitHubClient
	configManager  *ConfigManager
}

// NewUpdater creates a new updater instance
func NewUpdater(currentVersion string) (*Updater, error) {
	configManager, err := NewConfigManager()
	if err != nil {
		return nil, err
	}
	
	return &Updater{
		currentVersion: currentVersion,
		client:         NewGitHubClient("pyhub-apps", "pyhub-imagekit"),
		configManager:  configManager,
	}, nil
}

// CheckForUpdate checks if a new version is available
func (u *Updater) CheckForUpdate() (*Release, bool, error) {
	release, err := u.client.GetLatestRelease()
	if err != nil {
		return nil, false, err
	}
	
	isNewer := release.IsNewerThan(u.currentVersion)
	return release, isNewer, nil
}

// Update performs the self-update
func (u *Updater) Update(release *Release, force bool) error {
	// Get the appropriate asset for this platform
	asset := release.GetAssetForPlatform()
	if asset == nil {
		return fmt.Errorf("no binary available for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	
	fmt.Printf("다운로드 중... (%s)\n", asset.Name)
	
	// Download the new binary
	data, err := u.downloadBinary(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Create backup
	backupPath := execPath + ".backup"
	if err := u.createBackup(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	
	// Replace the binary
	if err := u.replaceBinary(execPath, data); err != nil {
		// Restore from backup on failure
		u.restoreBackup(backupPath, execPath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	
	// Remove backup
	os.Remove(backupPath)
	
	// Update config
	u.configManager.UpdateLastCheck(release.TagName)
	
	return nil
}

// downloadBinary downloads the binary from the given URL
func (u *Updater) downloadBinary(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

// createBackup creates a backup of the current binary
func (u *Updater) createBackup(srcPath, backupPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	
	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, src)
	return err
}

// replaceBinary replaces the current binary with new data
func (u *Updater) replaceBinary(execPath string, data []byte) error {
	// On Windows, we need to rename the old file first
	if runtime.GOOS == "windows" {
		oldPath := execPath + ".old"
		if err := os.Rename(execPath, oldPath); err != nil {
			return err
		}
		defer os.Remove(oldPath)
	}
	
	// Write new binary
	tempPath := execPath + ".new"
	if err := os.WriteFile(tempPath, data, 0755); err != nil {
		return err
	}
	
	// Move new binary to final location
	if err := os.Rename(tempPath, execPath); err != nil {
		return err
	}
	
	// Make it executable on Unix-like systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(execPath, 0755); err != nil {
			return err
		}
	}
	
	return nil
}

// restoreBackup restores the binary from backup
func (u *Updater) restoreBackup(backupPath, execPath string) error {
	// Remove failed update attempt
	os.Remove(execPath)
	
	// Restore backup
	return os.Rename(backupPath, execPath)
}

// ShowUpdateNotification displays an update notification if available
func (u *Updater) ShowUpdateNotification() {
	// Check if we should check for updates
	shouldCheck, err := u.configManager.ShouldCheckUpdate()
	if err != nil || !shouldCheck {
		return
	}
	
	// Check for updates
	release, hasUpdate, err := u.CheckForUpdate()
	if err != nil || !hasUpdate {
		return
	}
	
	// Update last check time
	u.configManager.UpdateLastCheck(u.currentVersion)
	
	// Show notification
	fmt.Println()
	fmt.Printf("✨ 새 버전이 있습니다! (현재: %s → 최신: %s)\n", u.currentVersion, release.TagName)
	fmt.Println("   업데이트: imagekit update")
	fmt.Printf("   변경사항: %s\n", release.HTMLURL)
}

// GetCurrentExecutablePath returns the path of the current executable
func GetCurrentExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	
	// Resolve any symlinks
	return filepath.EvalSymlinks(execPath)
}