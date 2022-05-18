package versions

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	build_counter "github.krafton.com/sbx/version-helper/pkg/modules/build-counter"
)

type SbxVersion struct {
	// base-version
	major uint
	minor uint
	patch uint

	// pre-release
	branch    string
	revision  uint
	commitSha string

	// flags
	exposeRevision bool

	originalBranch    string
	originalCommitSha string
	revisionCounter   build_counter.Counter
}

func NewSbxVersion(major uint, minor uint, patch uint) *SbxVersion {
	return &SbxVersion{major: major, minor: minor, patch: patch}
}

func NewDetailedSbxVersion(baseVersion *semver.Version, branch string, commitSha string, counter build_counter.Counter, exposeRevision bool) (*SbxVersion, error) {
	version := &SbxVersion{
		major:             uint(baseVersion.Major()),
		minor:             uint(baseVersion.Minor()),
		patch:             uint(baseVersion.Patch()),
		originalBranch:    branch,
		originalCommitSha: commitSha,
		revisionCounter:   counter,
		exposeRevision:    exposeRevision,
	}

	revision, err := counter.Get(context.TODO())
	if err != nil {
		return nil, err
	}
	version.revision = revision

	version.branch = MangleBranch(branch)
	version.commitSha = commitSha[0:7]

	return version, nil
}

func (v *SbxVersion) Major() uint {
	return v.major
}

func (v *SbxVersion) Minor() uint {
	return v.minor
}

func (v *SbxVersion) Patch() uint {
	return v.patch
}

func (v *SbxVersion) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)
	if preRelease := v.PreRelease(); preRelease != "" {
		fmt.Fprintf(&buf, "-%s", preRelease)
	}
	return buf.String()
}

func (v *SbxVersion) BaseVersion() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func (v *SbxVersion) Revision() uint {
	return v.revision
}

func (v *SbxVersion) PreRelease() string {
	if v.branch == "" {
		return ""
	}
	mangledBranch := MangleBranch(v.branch)
	if v.exposeRevision {
		if v.commitSha == "" {
			return fmt.Sprintf("%s.%d", mangledBranch, v.revision)
		}
		return fmt.Sprintf("%s.%d.%s", mangledBranch, v.revision, v.commitSha)
	} else {
		if v.commitSha == "" {
			return mangledBranch
		}
		return fmt.Sprintf("%s.%s", mangledBranch, v.commitSha)
	}
}

func (v *SbxVersion) Increase(ctx context.Context) error {
	revision, err := v.revisionCounter.Increase(ctx)
	if err != nil {
		return err
	}
	v.revision = revision
	return nil
}

func (v *SbxVersion) ToSemver() (*semver.Version, error) {
	return semver.StrictNewVersion(v.String())
}

func MangleBranch(branch string) string {
	return strings.ReplaceAll(strings.ToLower(branch), "/", "-")
}
