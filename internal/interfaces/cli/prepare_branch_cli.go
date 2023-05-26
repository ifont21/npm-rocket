/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"fmt"
	"os"

	"github.com/ifont21/pre-releaser-cli/internal/adapters"
	"github.com/ifont21/pre-releaser-cli/internal/domain"
	"github.com/spf13/cobra"
)

var repoLocalPath string

// bumpPackagesCmd represents the bumpPackages command
var prepareBranch = &cobra.Command{
	Use:   "prep-branch",
	Short: "Prepare releases for npm packages",
	Long:  `Create branch and commit the changes to the local repo branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if repoLocalPath == "" {
			repoLocalPath = dirPath
		}

		fmt.Println("Preparing branch for release...")
		fmt.Println("Repo path: ", repoLocalPath)

		fileRepository := adapters.NewFileRepository(repoLocalPath)
		config := adapters.NewConfig(fileRepository)
		gitChangesRepository := adapters.NewGitChanges(repoLocalPath)
		prGHRepository := adapters.NewPRGithubRepository(config)
		changesService := domain.NewAddChangesService(gitChangesRepository, prGHRepository)
		changesService.AddChanges()
	},
}

func init() {
	prepareBranch.Flags().StringVarP(&repoLocalPath, "local-repo", "r", "", "where the repo is located")

	rootCmd.AddCommand(prepareBranch)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpPackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpPackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
