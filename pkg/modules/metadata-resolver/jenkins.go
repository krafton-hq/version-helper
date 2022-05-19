package metadata_resolver

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-helper/pkg/modules/git"
	"go.uber.org/zap"
)

type JenkinsResolver struct {
}

const (
	jenkinsCheckEnv   = "JENKINS_URL"
	jenkinsCommitEnv  = "GIT_COMMIT"
	jenkinsRepoUrlEnv = "GIT_URL"
)

var jenkinsBranchEnvs = []string{"BRANCH_NAME", "GIT_BRANCH", "GIT_LOCAL_BRANCH"}

func (r *JenkinsResolver) String() string {
	return "Jenkins"
}

func (r *JenkinsResolver) CheckResolveTarget() bool {
	return os.Getenv(jenkinsCheckEnv) != ""
}

func (r *JenkinsResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	commitSha := os.Getenv(jenkinsCommitEnv)
	if commitSha == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: %s env not exists", jenkinsCommitEnv))
	}

	branch := git.NormalizeBranch(GetMultipleEnv(jenkinsBranchEnvs))
	if branch == "" {
		return nil, errors.New(fmt.Sprintf("EnvNotExists: none of %+v envs are exist", jenkinsBranchEnvs))
	}

	tag, err := git.GetTag()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("TagParseError: %s", err.Error()))
	}

	lastVersion, err := semver.NewVersion(tag)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SemverParseError: %s", err.Error()))
	}

	repo, err := git.GetRepositoryFromEnv(jenkinsRepoUrlEnv)
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

func GetMultipleEnv(envs []string) string {
	for _, env := range envs {
		if value, exists := os.LookupEnv(env); exists && value != "" {
			zap.S().Debugf("Found env %s=%s", env, value)
			return value
		}
	}
	return ""
}
