/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gogitty/internal/core"
	"gogitty/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash",
	Short: "A brief description of your command",
	Long:  `converts an existing file into a git object`,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		content, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		// create a new git blob object
		blob := core.GitBlob{}
		blob.Init()
		blob.BlobData = content

		cwd, _ := os.Getwd()
		repo, _ := utils.RepoFind(cwd, true)

		// write the object to the git object store
		hash, err := core.ObjectWrite(&blob, repo, true)

		if err != nil {
			panic(err)
		}
		// print the hash
		println(hash)

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
