/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gogitty/internal/core"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash",
	Short: "A brief description of your command",
	Long:  `converts an existing file into a git object`,
	Run: func(cmd *cobra.Command, args []string) {
		var fp string
		var stats fs.FileInfo
		var err error

		// check if custom path is provided or get current working directory
		if len(args) > 0 {
			fp = args[0]
			stats, err = os.Stat(fp) // Capture error from os.Stat
			if err != nil {
				// Handle the error appropriately, for example:
				fmt.Println("Error getting file stats:", err)
				return
			}
		} else {
			// If no argument is provided, you might want to set a default fp or exit
			fmt.Println("No file path provided.")
			return
		}

		data, err := os.ReadFile(fp)
		if err != nil {
			// Handle the error appropriately, for example:
			fmt.Println("Error reading file:", err)
			return
		}

		blobObject := core.Blob{
			Content:     data,
			ContentSize: stats.Size(), // Use stats here
		}

		currentRepoPath, err := core.RepoFind(fp, true)
		if err != nil {
			fmt.Println("Couldn't find git repo ")
		} else {
			blobObject.Serialize(currentRepoPath.GitDir)
		}

		// create directory os.MkdirAll(filepath.Dir(blobPath), os.ModePerm)
		os.MkdirAll(filepath.Dir(blobObject.BlobPath), os.ModePerm)
		f, err := os.Create(blobObject.BlobPath)
		if err != nil {
			fmt.Print("error occurred ", err.Error())
		}
		defer f.Close()
		f.Write(blobObject.Buffer.Bytes())
	},
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
