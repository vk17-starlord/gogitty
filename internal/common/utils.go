package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetRepoPath constructs a full path from a base repository path and additional subdirectories.
func GetRepoPath(repo string, path ...string) (string, error) {
	fullPath := filepath.Join(append([]string{repo}, path...)...)
	return fullPath, nil
}

// exists returns whether the given file or directory exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetRepoDir returns the directory path, creating it if it does not exist.
func GetRepoDir(repo string, mkdir bool, path ...string) (string, error) {
	dirpath, err := GetRepoPath(repo, path...)
	if err != nil {
		return "", fmt.Errorf("error constructing repo path: %w", err)
	}

	// Check if the directory exists
	if exists, err := exists(dirpath); err != nil {
		return "", fmt.Errorf("error checking directory: %w", err)
	} else if exists {
		// Return if the directory already exists
		return dirpath, nil
	}

	// Create the directory if it doesn't exist
	if mkdir {
		err := os.MkdirAll(dirpath, 0750)
		if err != nil {
			return "", fmt.Errorf("error creating directory: %w", err)
		}
	} else {
		return "", fmt.Errorf("directory does not exist and mkdir is set to false: %s", dirpath)
	}

	return dirpath, nil
}
