package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trungle-csv/rog-codegen/internal/codegen"
	"github.com/trungle-csv/rog-codegen/internal/util"
	"gopkg.in/yaml.v3"
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
	generateCmd.Flags().StringP("swaggerFile", "s", "swagger.yaml", "OpenAPI 3.0 spec file.")
	generateCmd.Flags().StringP("configFile", "c", "rog-codegen.yaml", "A YAML config file that controls rog-codegen behavior.")
	rootCmd.AddCommand(generateCmd)
}

func generateRun(cmd *cobra.Command) {
	flagSwaggerFile, err := cmd.Flags().GetString("swaggerFile")
	if err != nil {
		exitWithError("Please specify a path to OpenAPI 3.0 file.\n")
	}

	flagConfigFile, err := cmd.Flags().GetString("configFile")
	if err != nil {
		exitWithError("Please specify a path to configuration file.\n")
	}

	configFile, err := os.ReadFile(flagConfigFile)
	if err != nil {
		exitWithError("Failed to read config file: %s", err)
	}
	var config codegen.Configuration
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		exitWithError("Failed to parse config file: %s", err)
	}

	config = config.UpdateDefaultValues()
	if err := config.Validate(); err != nil {
		exitWithError("Configuration error: %s\n", err)
	}

	swagger, err := util.LoadSwagger(flagSwaggerFile)
	if err != nil {
		exitWithError("Failed to load swagger spec: %s\n", err)
	}

	codegenService := codegen.NewCodegenService(swagger, config)
	if err = codegenService.Generate(); err != nil {
		exitWithError("Failed to generate code: %s\n", err)
	}
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
