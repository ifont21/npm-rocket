/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"os"

	"github.com/ifont21/pre-releaser-cli/internal/adapters"
	"github.com/spf13/cobra"
)

var repoPath string
var preRelease bool

// bumpPackagesCmd represents the bumpPackages command
var bumpNPMPackages = &cobra.Command{
	Use:   "bump-npm-packages",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if repoPath == "" {
			repoPath = dirPath
		}

		preReleaseService := adapters.NewPreReleaserContainer(repoPath, os.Getenv("OPENAI_TOKEN"), preRelease)
		preReleaseService.PreReleasePackages()
	},
}

func init() {
	bumpNPMPackages.Flags().StringVarP(&repoPath, "local-repo", "r", "", "where the repo is located")
	bumpNPMPackages.Flags().BoolVarP(&preRelease, "pre-release", "p", false, "Whether to pre-release or not")

	rootCmd.AddCommand(bumpNPMPackages)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpPackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpPackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
