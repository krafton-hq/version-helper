package metadata_resolver

import (
	"os"
)

type TeamcityResolver struct {
}

const (
	teamcityCheckEnv = "TEAMCITY_VERSION"
)

func (r *TeamcityResolver) String() string {
	return "Teamcity"
}

func (r *TeamcityResolver) CheckResolveTarget() bool {
	return os.Getenv(teamcityCheckEnv) != ""
}

func (r *TeamcityResolver) ResolveBuildMetadata() (*BuildMetadata, error) {
	resolver := &LocalResolver{}
	return resolver.ResolveBuildMetadata()
}
