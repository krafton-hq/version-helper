package metadata_resolver

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-maker/pkg/modules/git"
)

type GhaResolver struct {
}

const (
	ghaCheckEnv  = "GITHUB_ACTION"
	ghaCommitEnv = "GITHUB_SHA"
)

var ghaBranchEnvs = []string{"GITHUB_BASE_REF", "GITHUB_REF_NAME"}

func (r *GhaResolver) String() string {
	return "GithubActions"
}

func (r *GhaResolver) CheckResolveTarget() bool {
	_, exists := os.LookupEnv(ghaCheckEnv)
	return exists
}

func (r *GhaResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	commitSha := os.Getenv(ghaCommitEnv)
	if commitSha == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", ghaCommitEnv))
	}

	branch := git.NormalizeBranch(GetMultipleEnv(ghaBranchEnvs))
	if branch == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: none of %+v envs are exist", ghaBranchEnvs))
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
