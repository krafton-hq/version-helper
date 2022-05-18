package build_counter

import "context"

type Counter interface {
	String() string
	Increase(ctx context.Context) (uint, error)
	Get(ctx context.Context) (uint, error)
}
