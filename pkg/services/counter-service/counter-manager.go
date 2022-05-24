package counter_service

import (
	"fmt"

	build_counter "github.com/krafton-hq/version-helper/pkg/modules/build-counter"
	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	"github.com/thediveo/enumflag"
)

type CounterFlag enumflag.Flag

const (
	//None
	Local CounterFlag = iota
	Network
)

var CounterFlags = map[CounterFlag][]string{
	//None:    {},
	Local:   {"local"},
	Network: {"network"},
}

const foxNamespace = "SBX-VERSION"

type Option struct {
	// Required
	Flag    CounterFlag
	Project string

	// Optional
	FoxAddr    string
	FoxDialTls bool
	LocalPath  string
}

func NewCounter(option *Option) (build_counter.Counter, error) {
	switch option.Flag {
	//case None:
	//	return build_counter.NewLocalCounter(option.LocalPath, option.Project)
	case Local:
		return build_counter.NewLocalCounter(option.LocalPath, option.Project)
	case Network:
		foxClient, err := fox_utils.NewClient(&fox_utils.Option{
			FoxAddr:    option.FoxAddr,
			FoxDialTls: option.FoxDialTls,
		})
		if err != nil {
			return nil, err
		}
		return build_counter.NewFoxCounter(option.Project, foxNamespace, foxClient), nil
	default:
		return nil, fmt.Errorf("UnknownCounterFlag, %v", option.Flag)
	}
}
