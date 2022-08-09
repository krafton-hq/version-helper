package cmd

import (
	"context"

	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newVersionUploadCommand() *cobra.Command {
	var (
		objectFile string
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Version Object 서버로 업로드",
	}

	cmd.Flags().StringVar(&objectFile, "file", "version.yaml", "Version Object File Path")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		obj, err := version_object_service.LoadVersionObj(objectFile)
		if err != nil {
			zap.S().Infof("Load Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		err = version_object_service.UploadVersionObj(ctx, obj)
		if err != nil {
			zap.S().Infof("Upload Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
	}
	return cmd
}

func init() {
	versionCmd.AddCommand(newVersionUploadCommand())
}
