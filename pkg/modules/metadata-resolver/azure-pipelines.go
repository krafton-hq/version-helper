package metadata_resolver

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-helper/pkg/modules/git"
)

type AzpResolver struct {
}

const (
	azpCheckEnv   = "SYSTEM_COLLECTIONURI"
	azpBranchEnv  = "BUILD_SOURCEBRANCH"
	azpCommitEnv  = "BUILD_SOURCEVERSION"
	azpRepoUrlEnv = "BUILD_REPOSITORY_URI"
)

func (r *AzpResolver) String() string {
	return "AzurePipelines"
}

func (r *AzpResolver) CheckResolveTarget() bool {
	_, exists := os.LookupEnv(azpCheckEnv)
	return exists
}

func (r *AzpResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	commitSha := os.Getenv(azpCommitEnv)
	if commitSha == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", azpCommitEnv))
	}

	branch := os.Getenv(azpBranchEnv)
	if branch == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: none of %+v envs are exist", azpBranchEnv))
	}

	tag, err := git.GetTag()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("TagParseError: %s", err.Error()))
	}

	lastVersion, err := semver.NewVersion(tag)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SemverParseError: %s", err.Error()))
	}

	repo, err := git.GetRepositoryFromEnv(azpRepoUrlEnv)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("RepositoryParseError: %s", err.Error()))
	}

	meta := &BuildMetadata{
		Branch:      branch,
		CommitSha:   commitSha,
		LastVersion: lastVersion,
		Repository:  repo,
	}
	return meta, nil
}
