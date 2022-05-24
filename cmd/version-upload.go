package cmd

import (
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
		Use:   "append",
		Short: "Version Object에 Artifact 추가",
	}

	cmd.Flags().StringVar(&uploadFoxAddr, "upload-fox-addr", "", "Network Counter Server gRPC Address")
	cmd.Flags().BoolVar(&uploadFoxTls, "upload-fox-secure", true, "Tls Flag to connect Network Counter Server")
	cmd.Flags().StringVar(&objectFile, "file", "version.yaml", "Version Object File Path")
	cmd.Flags().StringVarP(&conflictResolvePolicy, "conflict-resolve-policy", "c", "merge", "Resolve Policy, When Upload Conflicted")
	cmd.Flags().UintVarP(&conflictRetry, "conflict-retry", "r", 5, "Mac Upload Retry Count")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		_, err := version_object_service.LoadVersionObj(objectFile)
		if err != nil {
			zap.S().Infof("Load Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

	}
	return cmd
}

func init() {
	//versionCmd.AddCommand(newVersionUploadCommand())
}
