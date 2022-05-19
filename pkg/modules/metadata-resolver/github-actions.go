package metadata_resolver

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-helper/pkg/modules/git"
)

type GhaResolver struct {
}

const (
	ghaCheckEnv     = "GITHUB_ACTION"
	ghaCommitEnv    = "GITHUB_SHA"
	ghaRepoEnv      = "GITHUB_REPOSITORY"
	GhaRepoOwnerEnv = "GITHUB_REPOSITORY_OWNER"
)

var ghaBranchEnvs = []string{"GITHUB_BASE_REF", "GITHUB_REF_NAME", "GITHUB_REF"}

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

	// Expected: <owner>/<repo>
	githubRepo, exist := os.LookupEnv(ghaRepoEnv)
	if !exist {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", ghaRepoEnv))
	}

	// Expected: <owner>
	githubOwner, exist := os.LookupEnv(GhaRepoOwnerEnv)
	if !exist {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", GhaRepoOwnerEnv))
	}

	// Expected: <repo>
	repo := strings.TrimPrefix(githubRepo, githubOwner+"/")

	meta := &BuildMetadata{
		Branch:      branch,
		CommitSha:   commitSha,
		LastVersion: lastVersion,
		Repository:  repo,
	}
	return meta, nil
}
