package core

import (
	"errors"
	"fmt"
	"gogitty/internal/common"
	"gogitty/internal/config"
	"os"
	"path/filepath"
)

// Struct for storing Repository Object
type Repository struct {
	WorkTree string
	GitDir   string
}

// RepoFind searches for the repository root starting from the given path.
func RepoFind(path string, required bool) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("could not determine absolute path: %w", err)
	}

	for {
		// Check if .git directory exists in the current path
		gitDir := filepath.Join(absPath, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return &Repository{WorkTree: absPath, GitDir: gitDir}, nil
		}

		// Move up one directory level
		parent := filepath.Dir(absPath)
		if parent == absPath {
			// Reached the root
			if required {
				return nil, errors.New("no git directory found")
			}
			return nil, nil
		}

		// Set the new path to the parent
		absPath = parent
	}
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

	return nil
}
