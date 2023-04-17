/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"strings"

	"github.com/ifont21/pre-releaser-cli/internal"
	"github.com/spf13/cobra"
)

var basePath string
var libs string
var commits string

// bumpPackagesCmd represents the bumpPackages command
var bumpPackagesCmd = &cobra.Command{
	Use:   "bumpPackages",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		npmService := internal.NewNPMPackageService(internal.NewPkgJSONRepositoryImpl(), internal.NewChangeLogRepositoryImpl())
		libArray := strings.Split(libs, ",")
		err := npmService.BumpNPMPackagesAndChangelog(basePath, libArray, commits)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	bumpPackagesCmd.Flags().StringVarP(&basePath, "basePath", "b", "", "base path to the libraries")
	bumpPackagesCmd.Flags().StringVarP(&libs, "libs", "l", "", "libs to bump separated by comma")
	bumpPackagesCmd.Flags().StringVarP(&commits, "commits", "c", "", "commits to update the CHANGELOG.md file")

	rootCmd.AddCommand(bumpPackagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpPackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpPackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
