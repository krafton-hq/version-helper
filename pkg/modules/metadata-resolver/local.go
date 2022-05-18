package metadata_resolver

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.krafton.com/sbx/version-maker/pkg/modules/git"
)

type LocalResolver struct {
}

func (r *LocalResolver) String() string {
	return "Local"
}

func (r *LocalResolver) CheckResolveTarget() bool {
	return true
}

func (r *LocalResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	commitSha, err := git.GetCommit()
	if err != nil {
		return nil, err
	}

	branch, err := git.GetBranch()
	if err != nil {
		return nil, err
	}

	tag, err := git.GetTag()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("TagParseError: %s", err.Error()))
	}

	lastVersion, err := semver.NewVersion(tag)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("SemverParseError: %s", err.Error()))
	}

	repo, err := git.GetRepository()
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
