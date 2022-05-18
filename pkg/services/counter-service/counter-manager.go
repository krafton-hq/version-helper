package counter_service

import (
	"fmt"

	"github.com/thediveo/enumflag"
	build_counter "github.krafton.com/sbx/version-maker/pkg/modules/build-counter"
	"github.krafton.com/xtrm/fox/client/fox_grpc"
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
		foxConfig := fox_grpc.DefaultConfig()
		if option.FoxAddr != "" {
			foxConfig.GrpcEndpoint = option.FoxAddr
			foxConfig.WithTls = option.FoxDialTls
		}
		foxClient, err := fox_grpc.NewClient(foxConfig)
		if err != nil {
			return nil, err
		}

		return build_counter.NewFoxCounter(option.Project, foxNamespace, foxClient), nil
	default:
		return nil, fmt.Errorf("UnknownCounterFlag, %v", option.Flag)
	}
}
