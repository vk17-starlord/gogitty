/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gogitty/internal/core"
	"log"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command is used for initializing using folder for version control and tracking files",
	Long:  `Will create a .git hidden folder for version control`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create the repository
		repo := core.NewRepository()
		if repo == nil {
			log.Fatalf("Error creating repository")
		}

		repo.Init()
		fmt.Printf("Repository created at: %s\n", repo.Gitdir)

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
