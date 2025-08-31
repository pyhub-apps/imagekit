package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// Release represents a GitHub release
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
	HTMLURL     string  `json:"html_url"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
}

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	owner      string
	repo       string
	httpClient *http.Client
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(owner, repo string) *GitHubClient {
	return &GitHubClient{
		owner: owner,
		repo:  repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLatestRelease fetches the latest release from GitHub
func (gc *GitHubClient) GetLatestRelease() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", gc.owner, gc.repo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "imagekit-cli")
	
	resp, err := gc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}
	
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &release, nil
}

// GetAssetForPlatform returns the appropriate asset for the current platform
func (r *Release) GetAssetForPlatform() *Asset {
	platform := runtime.GOOS
	arch := runtime.GOARCH
	
	// Build the expected asset name
	var expectedName string
	switch platform {
	case "darwin":
		expectedName = fmt.Sprintf("imagekit-darwin-%s", arch)
	case "windows":
		expectedName = fmt.Sprintf("imagekit-windows-%s.exe", arch)
	case "linux":
		expectedName = fmt.Sprintf("imagekit-linux-%s", arch)
	default:
		return nil
	}
	
	// Find matching asset
	for _, asset := range r.Assets {
		if asset.Name == expectedName {
			return &asset
		}
	}
	
	return nil
}

// IsNewerThan checks if this release is newer than the given version
func (r *Release) IsNewerThan(currentVersion string) bool {
	// Remove 'v' prefix if present
	current := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(r.TagName, "v")
	
	// Simple string comparison for HeadVer format (head.yearweek.build)
	// This works because yearweek increases monotonically
	return latest > current
}

// GetChangesSummary returns a brief summary of changes
func (r *Release) GetChangesSummary() string {
	lines := strings.Split(r.Body, "\n")
	summary := []string{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			summary = append(summary, line)
			if len(summary) >= 3 {
				break // Limit to 3 items
			}
		}
	}
	
	if len(summary) == 0 {
		return "새로운 기능 및 버그 수정"
	}
	
	return strings.Join(summary, "\n")
}