package cmd

import (
	"context"
	"fmt"

	metadata_resolver "github.com/krafton-hq/version-helper/pkg/modules/metadata-resolver"
	meta_resolver_service "github.com/krafton-hq/version-helper/pkg/services/meta-resolver-service"
	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap"
)

func latestVersionGetCommand() *cobra.Command {
	var (
		ciHint meta_resolver_service.CiFlag

		namePrefix string
		repo       string
		branch     string
		namespace  string
	)

	cmd := &cobra.Command{
		Use:     "get {namePrefix}",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Short:   "최신 버전 조회 명령어",
		Example: "latestversion get MapDLC",
	}

	cmd.Flags().VarP(enumflag.New(&ciHint, "CI 이름", meta_resolver_service.CiFlags, enumflag.EnumCaseInsensitive),
		"ci-hint", "c", "CI 힌트")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "[REQUIRED] namespace name (ex: mapdlc-metadata)")
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "repository name")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "branch name")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		namePrefix = args[0]

		var metadata *metadata_resolver.BuildMetadata
		if repo == "" && branch == "" {
			resolver, err := GetMetaResolver(ciHint)
			if err != nil {
				zap.S().Infof("Get Resolver Failed, error: %s", err.Error())
				SetExitCode(ExitCodeError)
				return
			}
			zap.S().Debugf("Use Git Metadata Resolver: %s", resolver.String())

			if !resolver.CheckResolveTarget() {
				zap.S().Infof("Check Resolver Target Failed, Please check --ci-hint or Env")
				SetExitCode(ExitCodeError)
				return
			}

			metadata, err = resolver.ResolveBuildMetadata()
			if err != nil {
				zap.S().Infof("Resolve Build Metadata Failed, error: %s", err.Error())
				SetExitCode(ExitCodeError)
				return
			}
		}

		if repo == "" {
			repo = metadata.Repository
		}

		if repo == "" {
			zap.S().Infof("Repository is empty, Please check --repo or Env")
			SetExitCode(ExitCodeError)
			return
		}

		if branch == "" {
			branch = metadata.Branch
		}

		if branch == "" {
			zap.S().Infof("Branch is empty, Please check --branch or Env")
			SetExitCode(ExitCodeError)
			return
		}

		objName := namePrefix + "-" + repo + "-" + branch

		ctx := context.Background()

		result, err := version_object_service.GetLatestVersion(ctx, objName, namespace)
		if err != nil {
			zap.S().Infof("Get LatestVersion Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		fmt.Println(result.Status.VersionRef.Name)
	}
	return cmd
}

func init() {
	latestversionCmd.AddCommand(latestVersionGetCommand())
}
