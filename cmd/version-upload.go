package cmd

import (
	"context"

	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newVersionUploadCommand() *cobra.Command {
	var (
		objectFile            string
		conflictResolvePolicy string
		conflictRetry         uint

		uploadFoxAddr string
		uploadFoxTls  bool
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Version Object 서버로 업로드",
	}

	cmd.Flags().StringVar(&uploadFoxAddr, "upload-fox-addr", "", "Network Counter Server gRPC Address")
	cmd.Flags().BoolVar(&uploadFoxTls, "upload-fox-secure", true, "Tls Flag to connect Network Counter Server")
	cmd.Flags().StringVar(&objectFile, "file", "version.yaml", "Version Object File Path")
	cmd.Flags().StringVarP(&conflictResolvePolicy, "conflict-resolve-policy", "p", "merge", "Resolve Policy, When Upload Conflicted (one of merge or overwrite)")
	cmd.Flags().UintVarP(&conflictRetry, "conflict-retry", "r", 5, "Max Upload Retry Count")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		obj, err := version_object_service.LoadVersionObj(objectFile)
		if err != nil {
			zap.S().Infof("Load Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		err = version_object_service.UploadVersionObj(ctx, obj, &version_object_service.Option{
			ConflictResolvePolicy: conflictResolvePolicy,
			ConflictRetry:         int(conflictRetry),
			FoxAddr:               uploadFoxAddr,
			FoxDialTls:            uploadFoxTls,
		})
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
