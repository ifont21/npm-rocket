/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"os"
	"time"

	"github.com/ifont21/pre-releaser-cli/internal/adapters"
	"github.com/spf13/cobra"
)

var repoPath string
var preRelease bool
var dryRun bool
var noCommit bool

// bumpPackagesCmd represents the bumpPackages command
var prepareNPMPackagesRelease = &cobra.Command{
	Use:   "changelog-bumper",
	Short: "Prepare releases for npm packages",
	Long: `Read the local repo branch to find the packages that need to be released out of the commit messages.:

	1. Read the local repo branch to find the pre-releaser.yaml file and set up the configuration.
	2. Read the local repo branch to find the packages that need to be released out of the commit messages.
	3. Filter out commit messages that are not related to the packages that need to be released.
	4. Leveraging openai to bump the next version whether is a normal release or a pre-release.
	5. Leveraging openai to generate the changelog out of the commit messages. 
	
	Example:
	$ npm-rocket changelog-bumper -p`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if repoPath == "" {
			repoPath = dirPath
		}

		preReleaseService := adapters.NewPreReleaserContainer(repoPath, os.Getenv("OPENAI_TOKEN"), preRelease, dryRun)
		start := time.Now()
		preReleaseService.PreReleasePackages(preRelease, noCommit)
		elapsed := time.Since(start)
		cmd.Printf("Pre-release took %s", elapsed)
	},
}

func init() {
	prepareNPMPackagesRelease.Flags().StringVarP(&repoPath, "local-repo", "r", "", "where the repo is located")
	prepareNPMPackagesRelease.Flags().BoolVarP(&preRelease, "pre-release", "p", false, "Whether to pre-release or not")
	prepareNPMPackagesRelease.Flags().BoolVarP(&dryRun, "use-commit-file", "u", false, "Whether to dry run or not")
	prepareNPMPackagesRelease.Flags().BoolVarP(&noCommit, "no-commit", "n", false, "Whether to commit or not")

	rootCmd.AddCommand(prepareNPMPackagesRelease)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpPackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpPackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
