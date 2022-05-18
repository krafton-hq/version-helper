package metadata_resolver

import "github.com/Masterminds/semver/v3"

type Resolver interface {
	String() string
	CheckResolveTarget() bool
	ResolveBuildMetadata() (*BuildMetadata, error)
}

type BuildMetadata struct {
	Branch      string
	CommitSha   string
	Repository  string
	LastVersion *semver.Version
}
