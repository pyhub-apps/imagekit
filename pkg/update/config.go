package update

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the update configuration
type Config struct {
	LastUpdateCheck  time.Time `json:"lastUpdateCheck"`
	SkipUpdateCheck  bool      `json:"skipUpdateCheck"`
	LastKnownVersion string    `json:"lastKnownVersion"`
}

// ConfigManager handles configuration file operations
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new config manager
func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".pyhub", "imagekit")
	configPath := filepath.Join(configDir, "config.json")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	return &ConfigManager{
		configPath: configPath,
	}, nil
}

// Load reads the configuration from file
func (cm *ConfigManager) Load() (*Config, error) {
	config := &Config{}
	
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// Save writes the configuration to file
func (cm *ConfigManager) Save(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ShouldCheckUpdate determines if we should check for updates
func (cm *ConfigManager) ShouldCheckUpdate() (bool, error) {
	config, err := cm.Load()
	if err != nil {
		return false, err
	}
	
	if config.SkipUpdateCheck {
		return false, nil
	}
	
	// Check if 24 hours have passed since last check
	if time.Since(config.LastUpdateCheck) < 24*time.Hour {
		return false, nil
	}
	
	return true, nil
}

// UpdateLastCheck updates the last check timestamp
func (cm *ConfigManager) UpdateLastCheck(version string) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}
	
	config.LastUpdateCheck = time.Now()
	config.LastKnownVersion = version
	
	return cm.Save(config)
}