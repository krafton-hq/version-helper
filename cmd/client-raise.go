/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/krafton-hq/version-helper/pkg/modules/versions"
	counter_service "github.com/krafton-hq/version-helper/pkg/services/counter-service"
	meta_resolver_service "github.com/krafton-hq/version-helper/pkg/services/meta-resolver-service"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap"
)

func newClientRaiseCommand() *cobra.Command {
	var (
		tmplVersionFile string
		tmplHeaderFile  string
		genDir          string
		genVersionFile  string
		genHeaderFile   string

		project string

		ciHint meta_resolver_service.CiFlag

		counterFlag      counter_service.CounterFlag
		counterFoxAddr   string
		counterFoxTls    bool
		counterLocalPath string
	)

	cmd := &cobra.Command{
		Use:   "raise",
		Short: "Client Raise 버전 생성",
	}

	cmd.Flags().VarP(enumflag.New(&ciHint, "CI 이름", meta_resolver_service.CiFlags, enumflag.EnumCaseInsensitive),
		"ci-hint", "c", "CI 힌트")

	cmd.Flags().Var(enumflag.New(&counterFlag, "Counter 타입", counter_service.CounterFlags, enumflag.EnumCaseInsensitive),
		"counter", "Counter 타입 (local or network, default is local)")
	cmd.Flags().StringVar(&counterFoxAddr, "counter-fox-addr", "", "Network Counter Server gRPC Address")
	cmd.Flags().BoolVar(&counterFoxTls, "counter-fox-secure", true, "Tls Flag to connect Network Counter Server")
	cmd.Flags().StringVar(&counterLocalPath, "counter-local-path", "~/.versionhelper/db.json", "Local Counter DB File Path")

	cmd.Flags().StringVar(&project, "project", "client", "Project Name")

	cmd.Flags().StringVarP(&tmplVersionFile, "tmpl-version-file", "v", "embed:///client.yaml", "Template Version File Url (embded:///PATH, ./PATH or file:///PATH)")
	cmd.Flags().StringVarP(&tmplHeaderFile, "tmpl-header-file", "d", "embed:///GeneratedVersion.h", "Template Header File Url (embded:///PATH, ./PATH or file:///PATH)")
	cmd.Flags().StringVar(&genDir, "gen-dir", "", "Generated file output dir")
	cmd.Flags().StringVar(&genVersionFile, "gen-version-file", "version.yaml", "Version Metadata File Name (json or yaml)")
	cmd.Flags().StringVar(&genHeaderFile, "gen-header-file", "GeneratedVersion.h", "Header Metadata File Name (c header)")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		/**	Steps
		- Get Metadata
		- Create Counter
		- Create Version
		- Increase Count
		- Templating Output and Save
		*/

		// Get Metadata
		metadata, errored := clientGetMetadataStep(ciHint)
		if errored {
			return
		}

		// Create Counter
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

		// Create Version
		version, err := versions.NewDetailedSbxVersion(metadata.LastVersion, metadata.Branch, metadata.CommitSha, counter, true)
		if err != nil {
			zap.S().Infof("Create Sbx Version Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}

		// Increase Count
		err = version.Increase(context.TODO())
		if err != nil {
			zap.S().Infof("Increase Sbx Version Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		zap.S().Debugf("Version: %s", version.String())
		fmt.Println(version.String())

		// Templating Output and Save
		errored = clientTemplatingOutput(tmplVersionFile, genDir, genVersionFile, project, metadata, version)
		if errored {
			return
		}
		errored = clientTemplatingOutput(tmplHeaderFile, genDir, genHeaderFile, project, metadata, version)
		if errored {
			return
		}
	}

	return cmd
}

func init() {
	clientCmd.AddCommand(newClientRaiseCommand())
}
