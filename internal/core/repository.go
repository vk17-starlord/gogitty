package core

import (
	"fmt"
	"gogitty/internal/common"
	"gogitty/internal/config"
	"os"
)

// Struct for storing Repository Object
type Repository struct {
	WorkTree string
	GitDir   string
}

func (repo Repository) InitRepository() error {
	// Create necessary directories
	if _, err := common.GetRepoDir(repo.GitDir, true, "objects"); err != nil {
		return fmt.Errorf("error creating objects directory: %w", err)
	}
	if _, err := common.GetRepoDir(repo.GitDir, true, "refs", "heads"); err != nil {
		return fmt.Errorf("error creating refs/heads directory: %w", err)
	}
	if _, err := common.GetRepoDir(repo.GitDir, true, "refs", "tags"); err != nil {
		return fmt.Errorf("error creating refs/tags directory: %w", err)
	}

	// Create the HEAD file in the Git directory
	head, err := os.Create(repo.GitDir + "/HEAD")
	if err != nil {
		return fmt.Errorf("error creating HEAD file: %w", err)
	}
	defer head.Close()

	// Write HEAD file with default content
	if _, err := head.WriteString("ref: refs/heads/master\n"); err != nil {
		return fmt.Errorf("error writing to HEAD file: %w", err)
	}

	// Optionally write a config file
	_, err = config.New(repo.GitDir) // Assuming config.toml is in the git directory
	if err != nil {
		return fmt.Errorf("error initializing config: %w", err) // Return the error
	}

	return nil // Indicate success
}
