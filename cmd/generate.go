/*
Copyright © 2024 TRUNG LE <trunglq3007@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trungle-csv/rog-codegen/internal/util"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate API boilerplate",
	Run: func(cmd *cobra.Command, args []string) {
		generateRun(cmd)
	},
}

func init() {
	generateCmd.Flags().StringP("swagger", "s", "swagger.yaml", "Path to swagger file.")
	generateCmd.Flags().StringP("output", "o", "./", "Outut directory.")

	rootCmd.AddCommand(generateCmd)
}

func generateRun(cmd *cobra.Command) {
	filePath, err := cmd.Flags().GetString("swaggerFile")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid swagger file: %s\n", err)
		os.Exit(1)
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid output directory: %s\n", err)
		os.Exit(1)
	}

	_ = output
	swagger, err := util.LoadSwagger(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load swagger spec: %s\n", err)
		os.Exit(1)
	}

	_ = swagger
}
