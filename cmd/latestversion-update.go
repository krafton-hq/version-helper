package cmd

import (
	"context"
	"strings"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	metadata_resolver "github.com/krafton-hq/version-helper/pkg/modules/metadata-resolver"
	meta_resolver_service "github.com/krafton-hq/version-helper/pkg/services/meta-resolver-service"
	version_object_service "github.com/krafton-hq/version-helper/pkg/services/version-object-service"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func latestVersionUpdateCommand() *cobra.Command {
	var (
		ciHint meta_resolver_service.CiFlag

		namePrefix string
		versionRef string
		repo       string
		branch     string
		namespace  string
		labels     map[string]string
	)

	cmd := &cobra.Command{
		Use:     "update {namePrefix} {version}",
		Args:    cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Short:   "최신 버전 추가 또는 업데이트 명령어",
		Example: "latestversion update mapdlc abc123",
	}

	cmd.Flags().VarP(enumflag.New(&ciHint, "CI 이름", meta_resolver_service.CiFlags, enumflag.EnumCaseInsensitive),
		"ci-hint", "c", "CI 힌트")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "[REQUIRED] namespace name (ex: mapdlc-metadata)")
	cmd.Flags().StringVarP(&repo, "repo", "r", "", "repository name")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "branch name")
	cmd.Flags().StringToStringVarP(&labels, "label", "l", map[string]string{}, "latest version label (ex: hello=world)")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		namePrefix = args[0]
		versionRef = args[1]

		var metadata *metadata_resolver.BuildMetadata
		if repo == "" || branch == "" {
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

		objName := namePrefix + "-" + strings.ToLower(repo) + "-" + strings.ToLower(strings.ReplaceAll(branch, "/", "-"))

		ctx := context.Background()

		labels["branch"] = strings.ToLower(strings.ReplaceAll(branch, "/", "-"))
		labels["repository"] = repo

		lv := &redfoxV1alpha1.LatestVersion{
			TypeMeta: metav1.TypeMeta{
				Kind:       "LatestVersion",
				APIVersion: "metadata.sbx-central.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:   objName,
				Labels: labels,
			},
			Spec: redfoxV1alpha1.LatestVersionSpec{
				GitRef: redfoxV1alpha1.LatestVersionGitRef{
					Branch:     branch,
					Repository: repo,
				},
			},
			Status: redfoxV1alpha1.LatestVersionStatus{
				VersionRef: redfoxV1alpha1.LatestVersionVersionRef{
					Name: versionRef,
				},
			},
		}

		result, err := version_object_service.UploadLatestVersion(ctx, lv, namespace)
		if err != nil {
			zap.S().Infof("Upload LatestVersion Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		resultYaml, err := yaml.Marshal(result)
		if err != nil {
			zap.S().Infof("Marshal result LatestVersion Object Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		zap.S().Infof("Upload LatestVersion Object Success, name: %s", resultYaml)
	}
	return cmd
}

func init() {
	latestversionCmd.AddCommand(latestVersionUpdateCommand())
}
