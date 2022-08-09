package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version.yaml Object 파일 관리 명령어",
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
