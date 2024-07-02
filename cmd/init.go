/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trungle-csv/rog-codegen/internal/codegen"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a config file and OpenAPI 3.0 spec file.",
	Run: func(cmd *cobra.Command, args []string) {
		initRun()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun() {
	err := codegen.Initialize()
	if err != nil {
		exitWithError("Failed to create init files: %s", err)
	}

	fmt.Println("Generated .rog.yaml and openapi.yaml successfully!")
}
