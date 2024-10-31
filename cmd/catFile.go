/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"compress/zlib"
	"fmt"
	"gogitty/internal/common"
	"gogitty/internal/core"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "prints an existing git object to the standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		// check if custom path is provided or get current working directory
		var hash string
		if len(args) > 0 {
			hash = args[0]
		} else {
			return
		}
		dir, _ := os.Getwd()
		currentrepository, _ := core.RepoFind(dir, true)
		blobPath := fmt.Sprintf("%s/objects/%v/%v", currentrepository.GitDir, hash[0:2], hash[2:])

		data, err := os.Open(blobPath)
		common.CheckError(err)
		r, err := zlib.NewReader(io.Reader(data))
		if err != nil {
			panic(err)
		}
		s, _ := io.ReadAll(r)
		parts := strings.Split(string(s), "\x00")
		fmt.Println(parts[1])

	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
