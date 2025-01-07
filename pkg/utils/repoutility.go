package utils

import (
	"fmt"
	"gogitty/pkg/constants"
	"os"
	"path/filepath"
	"time"
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
	// Debugging output
	fmt.Println("Writing file to path:", path)
	fmt.Println("Content length:", len(content))
	fmt.Println("Content:", string(content))

	// Ensure the directory for the file exists.
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to ensure directory for file '%s': %w", path, err)
	}

	// Create or open the file (file will be overwritten if it already exists).
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", path, err)
	}
	defer file.Close()

	// Write the content to the file.
	n, err := file.Write(content)
	fmt.Printf("Wrote %d bytes to file.\n", n) // Debugging output
	if err != nil {
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

	gitFolderPath := filepath.Join(absPath, constants.GitFolder)
	if _, err := os.Stat(gitFolderPath); err != nil {
		parent := filepath.Join(absPath, "..")

		if parent == path {
			if required {
				return "", fmt.Errorf("no git directory")
			} else {
				return "", nil
			}
		}

	} else {
		return gitFolderPath, nil
	}

	return gitFolderPath, nil
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsExecutable(filePath string) (bool, error) {
	// Get file information
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	// Check if it's a regular file and not a directory
	if !info.Mode().IsRegular() {
		return false, nil
	}

	// Check if the file has execute permissions
	// Owner, group, or others can execute
	if info.Mode()&0111 != 0 {
		return true, nil
	}

	return false, nil
}

// GetTimestampAndOffset returns the current Unix timestamp and timezone offset in "+hhmm" format.
func GetTimestampAndOffset() (int64, string) {
	// Get the current time
	now := time.Now()

	// Get the Unix timestamp
	unixTimestamp := now.Unix()

	// Get the timezone offset in seconds
	_, offsetSeconds := now.Zone()
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60

	// Determine the sign and format the offset
	sign := "+"
	if offsetSeconds < 0 {
		sign = "-"
		offsetHours = -offsetHours
		offsetMinutes = -offsetMinutes
	}
	timezoneOffset := fmt.Sprintf("%s%02d%02d", sign, offsetHours, offsetMinutes)

	return unixTimestamp, timezoneOffset
}
