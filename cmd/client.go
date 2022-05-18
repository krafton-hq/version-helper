/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	metadata_resolver "github.krafton.com/sbx/version-maker/pkg/modules/metadata-resolver"
	"github.krafton.com/sbx/version-maker/pkg/modules/versions"
	generate_service "github.krafton.com/sbx/version-maker/pkg/services/generate-service"
	meta_resolver_service "github.krafton.com/sbx/version-maker/pkg/services/meta-resolver-service"
	"go.uber.org/zap"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client 버전 관리",
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func clientGetMetadataStep(ciHint meta_resolver_service.CiFlag) (*metadata_resolver.BuildMetadata, bool) {
	resolver, err := GetMetaResolver(ciHint)
	if err != nil {
		zap.S().Infof("Get Resolver Failed, error: %s", err.Error())
		SetExitCode(ExitCodeError)
		return nil, true
	}
	zap.S().Debugf("Use Git Metadata Resolver: %s", resolver.String())

	if !resolver.CheckResolveTarget() {
		zap.S().Infof("Check Resolver Target Failed, Please check --ci-hint or Env")
		SetExitCode(ExitCodeError)
		return nil, true
	}

	metadata, err := resolver.ResolveBuildMetadata()
	if err != nil {
		zap.S().Infof("Resolve Build Metadata Failed, error: %s", err.Error())
		SetExitCode(ExitCodeError)
		return nil, true
	}
	zap.S().Debugw("", "metadata", metadata)
	return metadata, false
}

func clientTemplatingOutput(tmplFile, genDir, genFile, project string, metadata *metadata_resolver.BuildMetadata, version versions.Version) bool {
	gen, err := generate_service.NewService(tmplFile, genDir, genFile)
	if err != nil {
		zap.S().Debugf("Init Generate Serviec Failed, error: %s", err.Error())
		SetExitCode(ExitCodeError)
		return true
	}

	err = gen.GenerateAndSave(version, metadata, project)
	if err != nil {
		zap.S().Debugf("Generate and Save Failed, error: %s", err.Error())
		SetExitCode(ExitCodeError)
		return true
	}
	return false
}
