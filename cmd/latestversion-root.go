package cmd

import (
	"github.com/spf13/cobra"
)

var latestversionCmd = &cobra.Command{
	Use:     "latestversion",
	Aliases: []string{"lv"},
	Short:   "latestversion 조회 및 업데이트 명령어",
}

func init() {
	rootCmd.AddCommand(latestversionCmd)
}
