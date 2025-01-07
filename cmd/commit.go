/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"gogitty/internal/core"
	"gogitty/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		commitMessage := args[0] // Get the commit message from arguments
		gitIndex := &core.GitIndex{}
		gitIndex.Init()
		// Step 1: Read the index file
		gitIndex.Deserialize()

		// Step 2: Read the HEAD file
		head := &core.GitHead{}
		head.Init()
		head.GetLatestCommitHash()

		// If the commit hash was found, create the parent commit object

		// Step 3: Create a new commit object
		commit := &core.GitCommit{}
		commit.Init("Author Name", "Committer Name", "Commit Message", head.Hash, "")

		fmt.Printf("parent hash is %s \n", head.Hash)
		path, _ := os.Getwd()
		repo, _ := utils.RepoFind(path, true)
		// Step 4: Build the tree object from the index
		treeHash := commit.BuildTreeFromIndex(gitIndex, repo)

		var buffer bytes.Buffer

		treeEntry := fmt.Sprintf("tree %s\n", treeHash)
		buffer.WriteString(treeEntry)
		timestamp, offset := utils.GetTimestampAndOffset()
		commitAuthor := fmt.Sprintf("author %s %d %s\n", "vinit khollam", timestamp, offset)
		buffer.WriteString(commitAuthor)
		committer := fmt.Sprintf("committer %s %d %s\n", "vinit khollam", timestamp, offset)
		buffer.WriteString(committer)
		buffer.WriteString(fmt.Sprintf("parent %s\n", head.Hash))
		buffer.WriteString(fmt.Sprintf("%s\n", commitMessage))

		commitContent := buffer.Bytes()

		header := fmt.Sprintf("commit %d\000\n", len(commitContent))
		fullContent := append([]byte(header), commitContent...)
		fmt.Printf(string(fullContent))

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
		head.AddLatestCommitToHead(hashHex)
		fmt.Printf(hashHex)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
