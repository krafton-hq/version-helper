package cmd

import (
	"os"

	"github.com/krafton-hq/version-helper/pkg/consts"
	log_helper "github.com/krafton-hq/version-helper/pkg/log-helper"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() {
	err := rootCmd.Execute()
	zap.S().Sync()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(flagExitCode)
}

var rootCmd = &cobra.Command{
	Use:     "versionhelper",
	Short:   "A brief description of your application",
	Version: consts.Version(),
}

func init() {
	rootCmd.SetVersionTemplate("version-helper: {{ .Version }}\n")

	var debug bool
	var jsonLog bool
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Flag to Print Debug Message")
	rootCmd.PersistentFlags().BoolVar(&jsonLog, "json-log", false, "Flag to Print Json Log Message")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		log_helper.Initialize(debug, jsonLog)
	}
}
