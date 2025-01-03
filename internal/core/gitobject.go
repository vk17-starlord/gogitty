package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gogitty/pkg/utils"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// GitObject interface
type GitObject interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	Init() error
	Format() string
}

func ObjectRead(repo string, sha string) {
	path := utils.RepoPath(repo, "objects", sha[:2], sha[2:])

	objectFile, err := os.ReadFile(path)
	utils.Check(err)

	// use zlib for reading content of the objectFILE
	data, err := zlib.NewReader(bytes.NewReader(objectFile))
	if err != nil {
		utils.Check(err)
	}

	rawContent, err := io.ReadAll(data)
	utils.Check(err)

	// Find the object type in the rawContent
	index := bytes.Index(rawContent, []byte(" "))

	// extract the type and size of the object
	y := bytes.Index(rawContent, []byte("\x00"))
	contentSize := string(rawContent[index:y])
	contentSize = strings.TrimSpace(contentSize)
	size, _ := strconv.Atoi(contentSize)

	// extract the rawContent of the object
	Actualcontent := rawContent[y+1:]
	if len(Actualcontent) == size {

		fmt.Print(string(Actualcontent))

	} else {

		fmt.Println("Content is not correct")
	}

	// close the reader
	defer data.Close()

}

func ObjectWrite(gitobject GitObject, repo string, write ...bool) (string, error) {

	// set default value of write to true
	writeFlag := true
	if len(write) > 0 {
		writeFlag = write[0]
	}

	// serialize the object
	data, _ := gitobject.Serialize()

	// read the file
	result := gitobject.Format() + " " + strconv.Itoa(len(data)) + "\x00" + string(data)

	hash := sha1.New()
	hash.Write([]byte(result))
	sha := fmt.Sprintf("%x", hash.Sum(nil))

	if writeFlag {
		Objectpath := utils.RepoPath(repo, "objects", sha[0:2])
		utils.EnsureDir(Objectpath)
		var buffer bytes.Buffer
		w := zlib.NewWriter(&buffer)
		w.Write([]byte(result))
		w.Close()
		err := os.WriteFile(Objectpath+"/"+sha[2:], buffer.Bytes(), 0644)

		if err != nil {
			return "", err
		}

	}

	return sha, nil

}

func ReadTree(repo string, sha string) {
	path := utils.RepoPath(repo, "objects", sha[:2], sha[2:])

	objectFile, err := os.ReadFile(path)
	utils.Check(err)

	// use zlib for reading content of the objectFILE
	data, err := zlib.NewReader(bytes.NewReader(objectFile))
	if err != nil {
		utils.Check(err)
	}
	defer data.Close()

	rawContent, err := io.ReadAll(data)

	// Find the object type in the rawContent
	SpaceIndex := bytes.Index(rawContent, []byte(" "))

	// extract the type and size of the object
	y := bytes.Index(rawContent, []byte("\x00"))
	contentSize := string(rawContent[SpaceIndex:y])
	contentSize = strings.TrimSpace(contentSize)

	// extract the rawContent of the object
	Actualcontent := rawContent[y+1:]

	// Step 2: Parse the decompressed data
	entries := []map[string]string{}
	treeData := Actualcontent
	index := 0

	// Step 3: Iterate over the decompressed data
	for index < len(treeData) {
		// Step 3a: Find file mode (until space character)
		modeEnd := bytes.IndexByte(treeData[index:], ' ')
		if modeEnd == -1 {
			fmt.Errorf("invalid tree data: mode not found")
			return
		}
		mode := string(treeData[index : index+modeEnd])
		index += modeEnd + 1

		// Step 3b: Find file name (until null byte)
		nameEnd := bytes.IndexByte(treeData[index:], '\x00')
		if nameEnd == -1 {
			fmt.Errorf("invalid tree data: filename not found")
		}
		filename := string(treeData[index : index+nameEnd])
		index += nameEnd + 1

		// Step 3c: Read the hash (next 20 bytes)
		if index+20 > len(treeData) {
			fmt.Errorf("invalid tree data: hash not found")
		}
		hash := treeData[index : index+20]
		index += 20

		// Step 3d: Convert the hash to hexadecimal
		hexHash := hex.EncodeToString(hash)

		// Step 4: Store the entry in the map
		entries = append(entries, map[string]string{
			"mode":     mode,
			"filename": filename,
			"hash":     hexHash,
		})
	}

	// Step 2: Print the decoded entries
	for _, entry := range entries {

		// Print the entry data
		fmt.Printf("Mode: %s, Filename: %s, Hash: %s\n", entry["mode"], entry["filename"], entry["hash"])
	}

}

func WriteTree() {
	// Step 1: Get the current working directory
	cwd, _ := os.Getwd()

	rootTree := readFiles(cwd, cwd)
	repo, _ := utils.RepoFind(cwd, true)
	_, treeHash, err := rootTree.Serialize(repo)
	if err != nil {
		fmt.Errorf("error serializing subtree: %s", err)
	} else {
		fmt.Print(treeHash)
	}
}

// readFiles reads files and directories recursively and adds them to the tree
func readFiles(path string, parent string) *GitTree {

	// List of filenames or patterns to ignore
	var ignoredFiles = []string{
		".git",         // Git directory
		".DS_Store",    // macOS system file
		"node_modules", // Common directory to ignore
		"*.tmp",        // Temporary files
		"*.log",        // Log files
		".jit",
	}

	tree := &GitTree{}
	tree.Init()

	// Step 1: Read the directory
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Errorf("error reading directory: %s", err)
	}

	// Step 2: Iterate over the files and directories
	for _, file := range files {
		// check if file is directory or normal file

		if slices.Contains(ignoredFiles, file.Name()) {
			continue
		}
		if file.IsDir() {
			// Step 2a: Add the directory to the tree
			subTree := readFiles(path+"/"+file.Name(), parent)

			repo, _ := utils.RepoFind(parent, true)

			_, treeHash, err := subTree.Serialize(repo)
			if err != nil {
				fmt.Errorf("error serializing subtree: %s", err)
			}

			tree.AddEntry("040000", treeHash, file.Name(), "tree")

			continue

		} else {
			// Step 2b: Add the file to the tree
			// Step 2b.1: Read the file content
			content, err := os.ReadFile(path + "/" + file.Name())
			if err != nil {
				fmt.Errorf("error reading file: %s", err)
			}
			// Step 2b.2: Create a new blob object
			blob := &GitBlob{}
			blob.Init()
			blob.BlobData = content
			// Step 2b.3: Write the blob object to the object store
			repo, _ := utils.RepoFind(parent, true)
			hash, err := ObjectWrite(blob, repo, true)
			if err != nil {
				fmt.Errorf("error writing blob: %s", err)
			}

			tree.AddEntry("100644", hash, file.Name(), "blob")
		}
	}
	return tree
}
