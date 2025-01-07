package core

import (
	"fmt"
	"gogitty/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

// GitHead represents the current HEAD of a Git repository
type GitHead struct {
	Ref  string
	Hash string
}

// Init initializes the GitHead object
func (gh *GitHead) Init() {
	gh.Ref = ""
	gh.Hash = ""
}

func (gh *GitHead) AddLatestCommitToHead(hash string) {
	// Get the current working directory
	cwd, _ := os.Getwd()

	// Find the repository
	repo, err := utils.RepoFind(cwd, true)
	if err != nil {
		fmt.Println("Error finding repository:", err)
		return
	}

	// Open the .git/HEAD file
	headFilePath := filepath.Join(repo, "HEAD")
	file, err := os.Open(headFilePath)
	if err != nil {
		fmt.Println("Error opening HEAD file:", err)
		return
	}

	defer file.Close()

	// Read the HEAD file content
	var refLine string
	_, err = fmt.Fscanf(file, "ref: %s\n", &refLine)
	if err != nil {
		fmt.Println("Error reading HEAD file:", err)
		return
	}

	if strings.HasPrefix(refLine, "refs/heads/") {
		branchName := strings.TrimPrefix(refLine, "refs/heads/")
		// Look for the commit hash under the refs directory
		refFilePath := filepath.Join(repo, "refs", "heads", branchName)
		err := os.WriteFile(refFilePath, []byte(hash), 0644)
		if err != nil {
			fmt.Println("Error writing commit hash to head file:", err)
			return
		}
	}

}

// GetLatestCommitHash fetches the latest commit hash by reading the HEAD file and checking the refs directory
func (gh *GitHead) GetLatestCommitHash() (string, bool) {
	// Get the current working directory
	cwd, _ := os.Getwd()

	// Find the repository
	repo, err := utils.RepoFind(cwd, true)
	if err != nil {
		fmt.Println("Error finding repository:", err)
		return "", false
	}

	// Open the .git/HEAD file
	headFilePath := filepath.Join(repo, "HEAD")
	file, err := os.Open(headFilePath)
	if err != nil {
		fmt.Println("Error opening HEAD file:", err)
		return "", false
	}
	defer file.Close()

	// Read the HEAD file content
	var refLine string
	_, err = fmt.Fscanf(file, "ref: %s\n", &refLine)
	if err != nil {
		fmt.Println("Error reading HEAD file:", err)
		return "", false
	}

	if strings.HasPrefix(refLine, "refs/heads/") {
		branchName := strings.TrimPrefix(refLine, "refs/heads/")
		// Look for the commit hash under the refs directory
		refFilePath := filepath.Join(repo, "refs", "heads", branchName)
		commitFile, err := os.Open(refFilePath)
		if err != nil {
			fmt.Println("Error opening branch file:", err)
			return "", false
		}
		defer commitFile.Close()

		// Read the commit hash from the branch file
		var commitHash string
		_, err = fmt.Fscanf(commitFile, "%s\n", &commitHash)
		if err != nil {
			return "", false
		}

		// Store the ref and hash in GitHead
		gh.Ref = branchName
		gh.Hash = commitHash
		return commitHash, true
	}

	return "", false
}

// PrintHead prints the current HEAD's ref and commit hash
func (gh *GitHead) PrintHead() {
	if gh.Ref != "" && gh.Hash != "" {
		fmt.Printf("Current HEAD: %s\nLatest commit hash: %s\n", gh.Ref, gh.Hash)
	} else {
		fmt.Println("No HEAD found or invalid HEAD file.")
	}
}
