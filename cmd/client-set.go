/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	build_counter "github.com/krafton-hq/version-helper/pkg/modules/build-counter"
	"github.com/krafton-hq/version-helper/pkg/modules/versions"
	meta_resolver_service "github.com/krafton-hq/version-helper/pkg/services/meta-resolver-service"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"go.uber.org/zap"
)

func newClientSetCommand() *cobra.Command {
	var (
		tmplFile string
		genDir   string
		genFile  string

		project string

		ciHint meta_resolver_service.CiFlag

		baseVersion *semver.Version
		count       uint
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Client Set 버전 생성",
		Args: cobra.MatchAll(
			cobra.ExactArgs(2),
			func(cmd *cobra.Command, args []string) error {
				// Set BaseVersion
				version, err := semver.StrictNewVersion(args[0])
				if err != nil {
					return err
				}
				if version.Prerelease() != "" || version.Metadata() != "" {
					return fmt.Errorf("UnexpectedVersion: Version should have only baseVersion")
				}
				baseVersion = version

				// Set Count
				if i, err := strconv.ParseUint(args[1], 10, 64); err != nil {
					return err
				} else {
					count = uint(i)
				}
				return nil
			}),
		Example: "versionhelper client set 0.3.23 231",
	}

	cmd.Flags().VarP(enumflag.New(&ciHint, "CI 이름", meta_resolver_service.CiFlags, enumflag.EnumCaseInsensitive),
		"ci-hint", "c", "CI 힌트")

	cmd.Flags().StringVar(&project, "project", "client", "Project Name")

	cmd.Flags().StringVarP(&tmplFile, "tmpl-file", "t", "embed:///client.yaml", "Template File Url (embded:///PATH, ./PATH or file:///PATH)")
	cmd.Flags().StringVar(&genDir, "gen-dir", "", "Generated file output dir")
	cmd.Flags().StringVar(&genFile, "gen-file", "version.yaml", "Version Metadata File Name (json or yaml)")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		/**	Steps
		- Get Metadata
		- Create Static Counter
		- Create Version
		- Templating Output and Save
		*/

		// Get Metadata
		metadata, errored := clientGetMetadataStep(ciHint)
		if errored {
			return
		}

		// Create Counter
		counter := build_counter.NewMemoryCounter(count)
		zap.S().Debugf("Use Counter: %s", counter.String())

		// Create Version
		version, err := versions.NewDetailedSbxVersion(baseVersion, metadata.Branch, metadata.CommitSha, counter, true)
		if err != nil {
			zap.S().Infof("Create Sbx Version Failed, error: %s", err.Error())
			SetExitCode(ExitCodeError)
			return
		}
		zap.S().Debugf("Version: %s", version.String())
		fmt.Println(version.String())

		// Templating Output and Save
		errored = clientTemplatingOutput(tmplFile, genDir, genFile, project, metadata, version)
		if errored {
			return
		}
	}

	return cmd
}

func init() {
	clientCmd.AddCommand(newClientSetCommand())
}
