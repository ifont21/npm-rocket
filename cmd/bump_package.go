/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ifont21/pre-releaser-cli/internal"
	"github.com/spf13/cobra"
)

var filePath string
var bumpType string

// bumpLibCmd represents the bumpLib command
var bumpPackageCmd = &cobra.Command{
	Use:   "bumpPkg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		npmService := internal.NewNPMPackageService(internal.NewPkgJSONRepositoryImpl(), nil)
		version, err := npmService.BumpNPMPackage(filePath, bumpType)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Bumped version to %s", version)
	},
}

func init() {
	bumpPackageCmd.Flags().StringVarP(&filePath, "file", "f", "", "package.json file path to bump the library")
	bumpPackageCmd.Flags().StringVarP(&bumpType, "type", "t", "", "bump type to update package.json version (major, minor, patch)")

	rootCmd.AddCommand(bumpPackageCmd)
}
