package core

import (
	"fmt"
	"gogitty/pkg/constants"
	"gogitty/pkg/utils"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	RepositoryFormatVersion int  `mapstructure:"repositoryformatversion"`
	FileMode                bool `mapstructure:"filemode"`
	Bare                    bool `mapstructure:"bare"`
}

type Repository struct {
	Worktree string
	Gitdir   string
	Config   string
}

func NewRepository() *Repository {
	//initialize the repository and .git should be take from constants package
	currentPath, err := os.Getwd()
	if err != nil {
		// handle the error appropriately
		return nil
	}
	return &Repository{Worktree: currentPath, Gitdir: currentPath + "/" + constants.GitFolder, Config: currentPath + constants.GitFolder + "/config"}
}

func (r *Repository) Init() {
	// Directories to create
	directories := []string{
		r.Gitdir,
		r.Gitdir + "/objects",
		r.Gitdir + "/refs",
		r.Gitdir + "/refs/heads",
		r.Gitdir + "/refs/tags",
	}

	// Create directories in a loop
	for _, dir := range directories {
		if err := utils.EnsureDir(dir); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Files to create
	files := map[string][]byte{
		r.Gitdir + "/HEAD":              []byte("ref: refs/heads/master"),
		r.Gitdir + "/config.yaml":       nil, // Config will be written later
		r.Gitdir + "/index":             []byte("HEADER:\n  Version: 2\n  Number of entries: 2"),
		r.Gitdir + "/refs/heads/master": nil,
	}

	// Create files in a loop
	for filePath, content := range files {
		if err := utils.CreateFile(filePath); err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			return
		}
		if content != nil {
			if err := utils.WriteFile(filePath, content); err != nil {
				fmt.Printf("Error writing to file %s: %v\n", filePath, err)
				return
			}
		}
	}

	// Write default config to the config file
	CreateConfig(r.Gitdir + "/config.yaml")
}

func CreateConfig(configPath string) error {
	// Set the config file path and format
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Check if the config file already exists
	if err := viper.ReadInConfig(); err != nil {
		// If the file doesn't exist, create a default configuration
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			defaultConfig := &Config{
				RepositoryFormatVersion: 0,
				FileMode:                false,
				Bare:                    false,
			}

			// Set default config values
			viper.Set("repositoryformatversion", defaultConfig.RepositoryFormatVersion)
			viper.Set("filemode", defaultConfig.FileMode)
			viper.Set("bare", defaultConfig.Bare)

			// Write the default config to the specified file
			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to create config file at '%s': %w", configPath, err)
			}
			fmt.Println("Config file created successfully!")
		} else {
			// Other errors related to reading the config
			return fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		// If the file exists, read and unmarshal its values
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			return fmt.Errorf("error unmarshalling config: %w", err)
		}
		fmt.Printf("Config loaded successfully: %+v\n", config)
	}

	// Optionally update configuration values
	fmt.Println("Updating config values...")
	viper.Set("repositoryformatversion", 1)
	viper.Set("filemode", true)
	viper.Set("bare", true)

	// Save the updated configuration back to the file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save updated config: %w", err)
	}
	return nil
}
