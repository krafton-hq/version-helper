package metadata_resolver

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-helper/pkg/modules/git"
)

type TeamcityResolver struct {
}

const (
	teamcityCheckEnv  = "TEAMCITY_VERSION"
	teamcityBranchEnv = ""
	teamcityCommitEnv = ""
)

func (r *TeamcityResolver) String() string {
	return "Teamcity"
}

func (r *TeamcityResolver) CheckResolveTarget() bool {
	return os.Getenv(teamcityCheckEnv) != ""
}

func (r *TeamcityResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	commitSha := os.Getenv(teamcityCommitEnv)
	if commitSha == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", teamcityCommitEnv))
	}

	branch := os.Getenv(teamcityBranchEnv)
	if branch == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: none of %+v envs are exist", teamcityBranchEnv))
	}

	tag, err := git.GetTag()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("TagParseError: %s", err.Error()))
	}

	lastVersion, err := semver.NewVersion(tag)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SemverParseError: %s", err.Error()))
	}

	meta := &BuildMetadata{
		Branch:      branch,
		CommitSha:   commitSha,
		LastVersion: lastVersion,
	}
	return meta, nil
}
