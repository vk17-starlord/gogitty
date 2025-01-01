package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gogitty/pkg/utils"
	"os"
)

// TreeEntry represents an entry in a Git tree (either a file or a subdirectory)
type TreeEntry struct {
	FileMode   string
	ObjectID   string
	FileName   string
	ObjectType string
}

// GitTree represents a Git tree object
type GitTree struct {
	Entries []TreeEntry
}

// Init initializes the GitTree object
func (t *GitTree) Init() error {
	t.Entries = []TreeEntry{}
	return nil
}

// AddEntry adds a new entry (file or subdirectory) to the tree
func (t *GitTree) AddEntry(fileMode string, objectID string, fileName string, objectType string) {
	t.Entries = append(t.Entries, TreeEntry{
		FileMode:   fileMode,
		ObjectID:   objectID,
		FileName:   fileName,
		ObjectType: objectType,
	})
}

// PrintEntries prints the entries of the GitTree in a formatted way
func (t *GitTree) PrintEntries() {
	// Check if there are any entries to print
	if len(t.Entries) == 0 {
		fmt.Println("No entries in the tree.")
		return
	}

	// Iterate over all entries and print in a human-readable format
	for _, entry := range t.Entries {
		fmt.Printf("File Name: %-20s | File Mode: %-8s | Object Type: %s | Object ID: %s\n",
			entry.FileName, entry.FileMode, entry.ObjectType, entry.ObjectID)
	}
}

// Serialize serializes the GitTree object into a byte slice and returns the SHA-1 hash
func (t *GitTree) Serialize(repo string) ([]byte, string, error) {
	var buffer bytes.Buffer

	// Serialize each tree entry
	for _, entry := range t.Entries {
		// Format: <file_mode> <file_name>\0<object_id>
		entryLine := fmt.Sprintf("%s %s\000", entry.FileMode, entry.FileName)
		buffer.WriteString(entryLine)
		objectIDBytes, err := hex.DecodeString(entry.ObjectID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode object ID: %v", err)
		}
		buffer.Write(objectIDBytes) // Append the raw 20-byte object ID
	}

	// Prepare the tree object content
	treeContent := buffer.Bytes()

	// Add the header: "tree <size>\0"
	header := fmt.Sprintf("tree %d\000", len(treeContent))
	fullContent := append([]byte(header), treeContent...)
	hash := sha1.New()
	hash.Write([]byte(fullContent))
	hashHex := fmt.Sprintf("%x", hash.Sum(nil))

	// write this tree files in object database
	Objectpath := utils.RepoPath(repo, "objects", hashHex[0:2])
	utils.EnsureDir(Objectpath)
	var TreeBuffer bytes.Buffer

	w := zlib.NewWriter(&TreeBuffer)
	w.Write(fullContent)
	w.Close()

	os.WriteFile(Objectpath+"/"+hashHex[2:], TreeBuffer.Bytes(), 0644)

	return fullContent, hashHex, nil
}

// Deserialize deserializes the byte slice into the GitTree object
func (t *GitTree) Deserialize(data []byte) error {

	return nil
}

// Format returns the object format, which is "tree" for GitTree
func (t *GitTree) Format() string {
	return "tree"
}
