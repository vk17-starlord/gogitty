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

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "prints an existing git object to the standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		// accept hash as argument
		hash := args[0]
		cwd, _ := os.Getwd()
		repo, _ := utils.RepoFind(cwd, true)
		// get the object
		core.ObjectRead(repo, hash)

	},
}

func init() {
	// Define the "type" flag with available choices
	catFileCmd.Flags().String("type", "", "The type of git object (blob, commit, tree, tag)")

	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
