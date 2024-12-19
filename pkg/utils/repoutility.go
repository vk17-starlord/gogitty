package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// RepoPath constructs a path within the repository by joining the base repo path with additional elements.
func RepoPath(repo string, paths ...string) string {
	elements := append([]string{repo}, paths...)
	return filepath.Join(elements...)
}

// EnsureDir ensures the directory exists, creating it if necessary.
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", path, err)
		}
	}
	return nil
}

// WriteFile writes content to a file at the specified path, creating any necessary directories.
func WriteFile(path string, content []byte) error {
	// Ensure the directory for the file exists.
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to ensure directory for file '%s': %w", path, err)
	}

	// Create and write to the file.
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", path, err)
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", path, err)
	}
	return nil
}

// CreateFile creates an empty file at the specified path, creating any necessary directories.
func CreateFile(path string) error {
	// Ensure the directory for the file exists.
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to ensure directory for file '%s': %w", path, err)
	}

	// Create the file.
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", path, err)
	}
	defer file.Close()

	return nil
}

func RepoFind(path string, required bool) (string, error) {
	// convert path to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Print("Error converting path to absolute path")
		return "", err
	}
	if _, err := os.Stat(filepath.Join(absPath, ".git")); err != nil {
		parent := filepath.Join(absPath, "..")

		if parent == path {
			// Bottom case
			// filepath.Join("/", "..") == "/":
			// If parent==path, then path is root.
			if required {
				return "", fmt.Errorf("no git directory")
			} else {
				return "", nil
			}
		}

	} else {
		return path, nil
	}

	return absPath, nil
}

