package versions

import "github.com/Masterminds/semver/v3"

type Version interface {
	String() string
	PreRelease() string
	BaseVersion() string
	Revision() uint
	ToSemver() (*semver.Version, error)
}
