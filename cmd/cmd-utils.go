package cmd

import (
	metadata_resolver "github.krafton.com/sbx/version-helper/pkg/modules/metadata-resolver"
	meta_resolver_service "github.krafton.com/sbx/version-helper/pkg/services/meta-resolver-service"
)

var flagExitCode int = 0

const ExitCodeOk = 0
const ExitCodeError = 1

func SetExitCode(exitCode int) {
	flagExitCode = exitCode
}

func GetMetaResolver(ciHint meta_resolver_service.CiFlag) (metadata_resolver.Resolver, error) {
	if ciHint == meta_resolver_service.None {
		return meta_resolver_service.SearchResolver(), nil
	} else {
		return meta_resolver_service.GetResolver(ciHint)
	}
}
