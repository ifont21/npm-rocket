/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/ifont21/pre-releaser-cli/internal"
	"github.com/spf13/cobra"
)

var changelogFilePath string

// Changelog text to be updated
var clText string

// updateChangelogCmd represents the updateChangelog command
var updateChangelogCmd = &cobra.Command{
	Use:   "updateChangelog",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		changelogService := internal.NewChangeLogPackageService(internal.NewChangeLogRepositoryImpl())
		newVersion := "## 3.2.1\n\n" +
			"- Testing changes\n\n"

		err := changelogService.UpdateChangelog(changelogFilePath, newVersion)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	updateChangelogCmd.Flags().StringVarP(&changelogFilePath, "file", "f", "", "CHANGELOG.md file path to be updated")
	updateChangelogCmd.Flags().StringVarP(&clText, "text", "t", "", "Changelog text to be updated")
	rootCmd.AddCommand(updateChangelogCmd)
}
