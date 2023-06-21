package cmd

import (
	"context"
	"fmt"

	"github.com/krafton-hq/version-helper/pkg/modules/versions"
	counter_service "github.com/krafton-hq/version-helper/pkg/services/counter-service"
	generate_service "github.com/krafton-hq/version-helper/pkg/services/generate-service"
	meta_resolver_service "github.com/krafton-hq/version-helper/pkg/services/meta-resolver-service"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap"
)

func newDevOpsServerCommand(teamName string) *cobra.Command {
	var (
		tmplFile string
		genDir   string
		genFile  string

		overrideProject string

		ciHint meta_resolver_service.CiFlag

		counterFlag      counter_service.CounterFlag
		counterFoxAddr   string
		counterFoxTls    bool
		counterLocalPath string
	)

	cmd := &cobra.Command{
		Use:   teamName,
		Short: fmt.Sprintf("%s 버전 생성", teamName),
	}

	cmd.Flags().VarP(enumflag.New(&ciHint, "CI 이름", meta_resolver_service.CiFlags, enumflag.EnumCaseInsensitive),
		"ci-hint", "c", "CI 힌트")

	cmd.Flags().Var(enumflag.New(&counterFlag, "Counter 타입", counter_service.CounterFlags, enumflag.EnumCaseInsensitive),
		"counter", "Counter 타입 (local or redfox, default is local)")
	cmd.Flags().StringVar(&counterLocalPath, "counter-local-path", "~/.versionhelper/db.json", "Local Counter DB File Path")

	cmd.Flags().StringVar(&overrideProject, "override-project", "", "Override Project Name (default use repository name)")

	cmd.Flags().StringVarP(&tmplFile, "tmpl-file", "t", "embed:///version.yaml", "Template File Url (embded:///PATH, ./PATH or file:///PATH)")
	cmd.Flags().StringVar(&genDir, "gen-dir", "", "Generated file output dir")
	cmd.Flags().StringVar(&genFile, "gen-file", "version.yaml", "Version Metadata File Name (json or yaml)")
	cmd.Run = func(cmd *cobra.Command, args []string) {
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

		metadata, err := resolver.ResolveBuildMetadata()
		if err != nil {
			zap.S().Infof("Resolve Build Metadata Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		zap.S().Debugw("", "metadata", metadata)

		project := metadata.Repository
		if overrideProject != "" {
			project = overrideProject
		}

		counter, err := counter_service.NewCounter(&counter_service.Option{
			Flag:       counterFlag,
			Project:    project,
			FoxAddr:    counterFoxAddr,
			FoxDialTls: counterFoxTls,
			LocalPath:  counterLocalPath,
		})
		if err != nil {
			zap.S().Infof("Create Counter Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		zap.S().Debugf("Use Counter: %s", counter.String())

		version, err := versions.NewDetailedSbxVersion(metadata.LastVersion, metadata.Branch, metadata.CommitSha, counter, false)
		if err != nil {
			zap.S().Infof("Create Sbx Version Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		err = version.Increase(context.TODO())
		if err != nil {
			zap.S().Infof("Increase Sbx Version Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		zap.S().Debugf("Version: %s", version.String())
		fmt.Println(version.String())

		gen, err := generate_service.NewService(tmplFile, genDir, genFile)
		if err != nil {
			zap.S().Debugf("Init Generate Serviec Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		err = gen.GenerateAndSave(version, metadata, project)
		if err != nil {
			zap.S().Debugf("Generate and Save Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
	}

	return cmd
}

func init() {
	rootCmd.AddCommand(newDevOpsServerCommand("devops"))
	rootCmd.AddCommand(newDevOpsServerCommand("server"))
	rootCmd.AddCommand(newDevOpsServerCommand("common"))
}
