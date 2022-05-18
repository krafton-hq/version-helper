package versions

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type Version interface {
	String() string
	PreRelease() string
	BaseVersion() string
	Revision() uint
	Major() uint
	Minor() uint
	Patch() uint
	Increase(ctx context.Context) error
	ToSemver() (*semver.Version, error)
}
