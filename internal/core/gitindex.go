package core

import (
	"bufio"
	"bytes"
	"fmt"
	"gogitty/pkg/constants"
	"gogitty/pkg/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GitIndex represents a Git index containing a map of entries
type GitIndex struct {
	Entries map[string]GitIndexEntry
}

// GitIndexEntry represents a single entry in the Git index
type GitIndexEntry struct {
	Mode string
	Hash string
}

// Init initializes the GitIndex object
func (gi *GitIndex) Init() {
	gi.Entries = make(map[string]GitIndexEntry)
}

// AddEntry adds a new entry to the GitIndex
func (gi *GitIndex) AddEntry(mode, hash, filename string) {
	entry := GitIndexEntry{
		Mode: mode,
		Hash: hash,
	}
	gi.Entries[filename] = entry
}

// PrintEntries prints all entries in the GitIndex
func (gi *GitIndex) PrintEntries() {
	for filename, entry := range gi.Entries {
		fmt.Println(entry.Mode, entry.Hash, filename)
	}
}

func (gi *GitIndex) Serialize() []byte {
	var buffer bytes.Buffer
	for filename, entry := range gi.Entries {
		buffer.WriteString(entry.Mode + " " + entry.Hash + " " + filename + "\n")
	}

	cwd, _ := os.Getwd()

	repo, _ := utils.RepoFind(cwd, true)

	// open the index file and add all entries to it
	file, err := os.OpenFile(repo+"/index", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	defer file.Close()

	// write the buffer to the index file
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}

	// close the index file
	err = file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		return nil
	}

	return buffer.Bytes()

}

func (gi *GitIndex) Deserialize() {

	wd, _ := os.Getwd()
	repo, _ := utils.RepoFind(wd, true)

	// Open the index file
	filePath := utils.RepoPath(repo, "index")
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by spaces (you can change this based on your format)
		parts := strings.Fields(line)

		if len(parts) >= 3 {
			// Extract file path and store it in the entries slice
			filePath := parts[2]
			gi.AddEntry(parts[0], parts[1], filePath)
		}
	}

	// Close the file
	err = file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}

}

func (gi *GitIndex) AddFileEntry(file fs.DirEntry, parent string) {

	var ignoredFiles = []string{
		".git",         // Git directory
		".DS_Store",    // macOS system file
		"node_modules", // Common directory to ignore
		"*.tmp",        // Temporary files
		"*.log",        // Log files
		".jit",
	}

	// Construct the full file path relative to the parent directory
	filePath := filepath.Join(parent, file.Name())

	// If it's a directory, recursively call AddFileEntry for its contents
	if file.IsDir() {
		// Skip directories that are ignored (e.g., .git)
		for _, ignoredDir := range ignoredFiles {
			if match, _ := filepath.Match(ignoredDir, file.Name()); match {
				return
			}
		}

		// Recurse into subdirectories
		files, err := os.ReadDir(filePath)
		if err != nil {
			log.Fatalf("Failed to read directory %s: %v", filePath, err)
		}

		// Process each file in the directory
		for _, subFile := range files {
			gi.AddFileEntry(subFile, filePath) // Pass full path for subdirectories
		}

		return
	}

	// Skip files that are in the ignored list
	for _, ignoredFile := range ignoredFiles {
		if match, _ := filepath.Match(ignoredFile, file.Name()); match {
			return
		}
	}

	// Step 1: Read file content
	content, err := os.ReadFile(filePath) // Use the full file path
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", filePath, err)
	}

	// Step 2: Create a new blob object
	blob := &GitBlob{}
	blob.Init()
	blob.BlobData = content

	// Step 2b: Write the blob object to the object store
	repopath, _ := os.Getwd()
	repoObj, err := utils.RepoFind(repopath, true)
	if err != nil {
		log.Fatalf("Error finding repo: %v", err)
	}

	hash, err := ObjectWrite(blob, repoObj, true)
	if err != nil {
		log.Fatalf("Error writing blob: %v", err)
	}

	rootPath := strings.Replace(repoObj, constants.GitFolder, "", -1)
	// Trim the root path from the absolute path
	relativePath := strings.TrimPrefix(filePath, rootPath)
	gi.AddEntry("100644", hash, relativePath)
}
