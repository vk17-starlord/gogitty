/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gogitty/internal/core"
	"os"

	"github.com/spf13/cobra"
)

var (
	repo core.Repository // Declare a package-level variable for the repo
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command is used for initializing using folder for version control and tracking files",
	Long:  `Will create a .git hidden folder for version control`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoPath string

		// check if custom path is provided or get current working directory
		if len(args) > 0 {
			repoPath = args[0]
		} else {
			var err error
			repoPath, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}
		}

		// pass file path and required configs
		repo := core.Repository{
			WorkTree: repoPath,
			GitDir:   fmt.Sprintf("%s/.git", repoPath),
		}

		// initialize repository
		repo.InitRepository()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
