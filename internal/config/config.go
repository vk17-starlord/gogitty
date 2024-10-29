package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the configuration structure.
type Config struct {
	Core struct {
		RepositoryFormatVersion string `mapstructure:"repositoryformatversion"`
		FileMode                bool   `mapstructure:"filemode"`
		Bare                    bool   `mapstructure:"bare"`
	} `mapstructure:"core"`
}

// New initializes a new configuration instance.
func New(repoDir string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(repoDir)

	// Set default values
	v.SetDefault("core.repositoryformatversion", "0")
	v.SetDefault("core.filemode", false)
	v.SetDefault("core.bare", false)

	// Check if the config file exists
	configFilePath := fmt.Sprintf("%s/config.toml", repoDir)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Create a new config file with default values
		if err := createDefaultConfig(configFilePath); err != nil {
			return nil, fmt.Errorf("error creating default config file: %w", err)
		}
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the config into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// createDefaultConfig creates a default config.toml file with predefined settings.
func createDefaultConfig(filePath string) error {
	defaultConfig := `[core]
repositoryformatversion = "0"
filemode = false
bare = false
`
	// Write the default configuration to the file
	return os.WriteFile(filePath, []byte(defaultConfig), 0644) // 0644 sets permissions for the file
}

// Save saves the configuration to a file.
func (c *Config) Save(repoDir string) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(repoDir)

	// Set the config values from the Config struct
	v.Set("core.repositoryformatversion", c.Core.RepositoryFormatVersion)
	v.Set("core.filemode", c.Core.FileMode)
	v.Set("core.bare", c.Core.Bare)

	// Save the config file
	if err := v.WriteConfigAs(repoDir + "/config.toml"); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
