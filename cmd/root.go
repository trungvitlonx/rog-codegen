/*
Copyright © 2024 TRUNG LE <trunglq3007@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rog-codegen",
	Short: "Generate Ruby on Rails API server boilerplate from OpenAPI 3 specifications",
	Long:  `Generate Ruby on Rails API server boilerplate from OpenAPI 3 specifications`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
