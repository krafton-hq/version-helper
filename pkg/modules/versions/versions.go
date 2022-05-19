package versions

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type Version interface {
	String() string
	BaseVersion() string
	PreRelease() string

	Major() uint
	Minor() uint
	Patch() uint
	Branch() string
	Revision() uint
	Commit() string

	Increase(ctx context.Context) error
	ToSemver() (*semver.Version, error)
}
