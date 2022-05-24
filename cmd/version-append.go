package cmd

import (
	"errors"

	version_object "github.com/krafton-hq/version-helper/pkg/modules/version-object"
	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newVersionAppendCommand() *cobra.Command {
	var (
		objectFile   string
		platform     string
		target       string
		artifactType string
		uri          string
		description  string
	)

	cmd := &cobra.Command{
		Use:   "append",
		Short: "Version Object에 Artifact 추가",
	}

	cmd.Flags().StringVar(&platform, "platform", "", "[Required] Artifact's Execute Format (ex: windows/amd64)")
	cmd.Flags().StringVar(&artifactType, "type", "", "[Required] Artifact's archived format (ex: zip, oci)")
	cmd.Flags().StringVar(&uri, "uri", "", "[Required] Artifact's Uri (ex: s3 or oci uri)")
	cmd.Flags().StringVar(&target, "target", "", "Artifact's Distribute Target (ex: dev, ship/appstore)")
	cmd.Flags().StringVar(&description, "description", "", "Artifact's Description")

	cmd.Flags().StringVar(&objectFile, "file", "version.yaml", "Version Object File Path")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if platform == "" || artifactType == "" || uri == "" {
			zap.S().Infof("The --platform, --type, --uri parameter should not be empty")
			return errors.New("invalid Parameters")
		}
		return nil
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		obj, err := version_object_service.LoadVersionObj(objectFile)
		if err != nil {
			zap.S().Infof("Load Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		artifact := &version_object.Artifact{
			Platform:     platform,
			Target:       target,
			ArtifactType: artifactType,
			Uri:          uri,
			Description:  description,
		}

		obj.Status.Artifact = append(obj.Status.Artifact, artifact)

		err = version_object_service.SaveVersionObj(obj, objectFile)
		if err != nil {
			zap.S().Infof("Save Version Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
	}

	return cmd
}

func init() {
	versionCmd.AddCommand(newVersionAppendCommand())
}
