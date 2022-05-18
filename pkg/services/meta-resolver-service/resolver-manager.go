package meta_resolver_service

import (
	"fmt"

	"github.com/thediveo/enumflag"
	metadata_resolver "github.krafton.com/sbx/version-helper/pkg/modules/metadata-resolver"
	"go.uber.org/zap"
)

type CiFlag enumflag.Flag

const (
	None CiFlag = iota
	Teamcity
	Jenkins
	GithubActions
	AzurePipelines
	Local
)

var CiFlags = map[CiFlag][]string{
	None:           {"none"},
	Teamcity:       {"teamcity", "tc"},
	Jenkins:        {"jenkins"},
	GithubActions:  {"github-actions", "actions", "github", "gha"},
	AzurePipelines: {"azure-pipelines", "azure-devops", "azp", "azdo"},
	Local:          {"local"},
}

func SearchResolver() metadata_resolver.Resolver {
	server := []metadata_resolver.Resolver{
		&metadata_resolver.AzpResolver{},
		&metadata_resolver.JenkinsResolver{},
		&metadata_resolver.TeamcityResolver{},
		&metadata_resolver.GhaResolver{},
	}
	for _, resolver := range server {
		if resolver.CheckResolveTarget() {
			return resolver
		} else {
			zap.S().Debugf("Check Resolver %s Status Failed, Skip this resolver", resolver.String())
		}
	}
	zap.S().Debugf("Check All CI Resolver Failed, Use local resolver")
	return &metadata_resolver.LocalResolver{}
}

func GetResolver(ci CiFlag) (metadata_resolver.Resolver, error) {
	switch ci {
	case Teamcity:
		return &metadata_resolver.TeamcityResolver{}, nil
	case Jenkins:
		return &metadata_resolver.JenkinsResolver{}, nil
	case GithubActions:
		return &metadata_resolver.GhaResolver{}, nil
	case AzurePipelines:
		return &metadata_resolver.AzpResolver{}, nil
	case Local:
		return &metadata_resolver.LocalResolver{}, nil
	default:
		return nil, fmt.Errorf("UnknownCiFlag, %v", ci)
	}
}
