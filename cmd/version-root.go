package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version Object 관리",
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
