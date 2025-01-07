package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gogitty/pkg/utils"
	"os"
	"sort"
	"strings"
	"time"
)

// GitCommit represents a Git commit object
type GitCommit struct {
	CommitHash string    // The commit's unique identifier (SHA-1 hash)
	Author     string    // Author name and email
	Committer  string    // Committer name and email
	Timestamp  time.Time // Timestamp of the commit
	Message    string    // The commit message
	ParentHash string    // Parent commit hash
	TreeHash   string    // The hash of the tree object associated with the commit
}

// Init initializes a new GitCommit with default values
func (commit *GitCommit) Init(author, committer, message, parentHash, treeHash string) {
	commit.Author = author
	commit.Committer = committer
	commit.Message = message
	commit.ParentHash = parentHash
	commit.TreeHash = treeHash
	commit.Timestamp = time.Now()
}

type FileNode struct {
	FileMode   string
	ObjectID   string
	FileName   string
	ObjectType string
	IsDir      bool
	Children   []*FileNode
}

type GitObjectTree struct {
	Name     string
	Hash     string
	Mode     string
	Children []*GitObjectTree
}

// Function to recursively insert a new file or directory into the tree
func insertGitObject(root *GitObjectTree, path string, hash, mode string) {
	// Split the path by "/"
	parts := strings.Split(path, "/")

	// If there are no parts left, this is a file
	if len(parts) == 0 {
		return
	}

	// Check if the current root is a directory (name is not empty)
	if len(parts) == 1 {
		// This is a file, so append it to the current directory
		root.Children = append(root.Children, &GitObjectTree{
			Name: parts[0],
			Hash: hash,
			Mode: mode,
		})

		// Sort the children alphabetically by their name
		sort.Slice(root.Children, func(i, j int) bool {
			return root.Children[i].Name < root.Children[j].Name
		})
		return
	}

	// Check if the current directory exists
	for _, child := range root.Children {
		// If we find a directory matching the first part of the path, recurse into it
		if child.Name == parts[0] {
			insertGitObject(child, strings.Join(parts[1:], "/"), hash, mode)
			// Sort the children of the directory alphabetically after insertion
			sort.Slice(child.Children, func(i, j int) bool {
				return child.Children[i].Name < child.Children[j].Name
			})
			return
		}
	}

	// If no directory was found, create a new one
	newDir := &GitObjectTree{
		Name: parts[0],
	}

	// Recursively insert the file or directory into the new directory
	insertGitObject(newDir, strings.Join(parts[1:], "/"), hash, mode)

	// Add the new directory to the root and sort the children alphabetically
	root.Children = append(root.Children, newDir)
	sort.Slice(root.Children, func(i, j int) bool {
		return root.Children[i].Name < root.Children[j].Name
	})
}

func (commit *GitCommit) BuildTreeFromIndex(gitIndex *GitIndex, repo string) string {
	root := &GitObjectTree{
		Name:     "root",
		Hash:     "",
		Mode:     "",
		Children: []*GitObjectTree{},
	}

	for filePath, entry := range gitIndex.Entries {
		insertGitObject(root, filePath, entry.Hash, entry.Mode)
	}

	// Print the tree structure recursively
	treeHash, _ := SerializeTree(root, repo)

	return treeHash
}

func SerializeTree(node *GitObjectTree, repo string) (string, error) {
	var buffer bytes.Buffer

	// Iterate over the children (files and directories)
	for _, child := range node.Children {
		if len(child.Children) > 0 {
			// Serialize the directory and get its tree hash
			treeHash, err := SerializeTree(child, repo) // Recursively serialize the directory
			if err != nil {
				return "", fmt.Errorf("failed to serialize directory %s: %v", child.Name, err)
			}

			// Convert tree hash to binary (20 bytes)
			treeHashBinary, err := hex.DecodeString(treeHash)
			if err != nil {
				return "", fmt.Errorf("invalid tree hash: %v", err)
			}

			// Write directory entry: <mode> <name>\0<binary hash>
			mode := "040000"
			if child.Mode != "" {
				mode = child.Mode
			}
			entry := fmt.Sprintf("%s %s", mode, child.Name)
			buffer.WriteString(entry)
			buffer.WriteByte(0) // Null byte separator
			buffer.Write(treeHashBinary)
		} else {
			// Convert blob hash to binary (20 bytes)
			blobHashBinary, err := hex.DecodeString(child.Hash)
			if err != nil {
				return "", fmt.Errorf("invalid blob hash: %v", err)
			}

			// Write file entry: <mode> <name>\0<binary hash>
			entry := fmt.Sprintf("%s %s", child.Mode, child.Name)
			buffer.WriteString(entry)
			buffer.WriteByte(0) // Null byte separator
			buffer.Write(blobHashBinary)
		}
	}

	// Prepare the tree content
	treeContent := buffer.Bytes()

	// Add the header: "tree <size>\0"
	header := fmt.Sprintf("tree %d\000", len(treeContent))
	fullContent := append([]byte(header), treeContent...)

	// Calculate the SHA-1 hash of the entire tree object
	hash := sha1.New()
	hash.Write(fullContent)
	hashHex := fmt.Sprintf("%x", hash.Sum(nil))

	// Construct the path for the object storage (objects/<prefix>/<hash>)
	objectPath := utils.RepoPath(repo, "objects", hashHex[:2])
	utils.EnsureDir(objectPath)

	// Write the tree object to a compressed file (zlib format)
	var treeBuffer bytes.Buffer
	w := zlib.NewWriter(&treeBuffer)
	w.Write(fullContent)
	w.Close()

	// Save the compressed tree object to disk
	err := os.WriteFile(objectPath+"/"+hashHex[2:], treeBuffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write tree object: %v", err)
	}

	return hashHex, nil
}
