package cmd

import (
	"errors"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newVersionAppendCommand() *cobra.Command {
	var (
		objectFile   string
		platform     string
		name         string
		artifactType string
		uri          string
		humanUri     string
		description  string
		labels       map[string]string
	)

	cmd := &cobra.Command{
		Use:   "append",
		Short: "Version Object에 Artifact 추가",
	}

	cmd.Flags().StringVar(&platform, "platform", "", "[Required if not label] Artifact's Execute Format (ex: windows/amd64)")
	cmd.Flags().StringVar(&artifactType, "type", "", "[Required if not label] Artifact's archived format (ex: zip, oci)")
	cmd.Flags().StringVar(&name, "name", "", "[Required if not label] Artifact's Name")
	cmd.Flags().StringVar(&uri, "uri", "", "[Required if not label] Artifact's Uri (ex: s3 or oci uri)")
	cmd.Flags().StringVar(&humanUri, "human-uri", "", " Artifact's Human Friendly Uri (ex: minio uri)")
	cmd.Flags().StringVar(&description, "description", "", "Artifact's Description")
	cmd.Flags().StringToStringVarP(&labels, "label", "l", map[string]string{}, "[Required if not platform, type, name, uri] Artifact's Label (ex: hello=world)")

	cmd.Flags().StringVar(&objectFile, "file", "version.yaml", "Version Object File Path")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if len(labels) == 0 && (platform == "" || artifactType == "" || uri == "" || name == "") {
			zap.S().Infof("The --label or --platform, --type, --uri --name parameter should not be empty")
			return errors.New("invalid Parameters")
		} else if len(labels) > 0 && (platform != "" || artifactType != "" || uri != "" || name != "") {
			zap.S().Infof("The --label or --platform, --type, --uri --name parameter should not be used together")
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

		if len(labels) > 0 {
			for k, v := range labels {
				obj.Labels[k] = v
			}
		} else {
			obj.Status.Artifacts = append(obj.Status.Artifacts, redfoxV1alpha1.VersionStatusArtifact{
				Name:             name,
				Type:             artifactType,
				Uri:              uri,
				Platform:         platform,
				HumanFriendlyUri: humanUri,
				Description:      description,
			})
		}

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
