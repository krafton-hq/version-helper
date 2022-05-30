package fox_utils

import (
	"fmt"

	log_helper "github.com/krafton-hq/version-helper/pkg/log-helper"
	"github.krafton.com/xtrm/fox/client/fox_grpc"
	"github.krafton.com/xtrm/fox/core/generated/protos"
)

func NewClient(option *Option) (*fox_grpc.FoxClient, error) {
	foxConfig := fox_grpc.DefaultConfig()
	if option.FoxAddr != "" {
		foxConfig.GrpcEndpoint = option.FoxAddr
		foxConfig.WithTls = option.FoxDialTls
	}
	foxConfig.ClientInterceptors = log_helper.GetUnaryClientInterceptors()

	foxClient, err := fox_grpc.NewClient(foxConfig)
	if err != nil {
		return nil, err
	}
	return foxClient, nil
}

type Option struct {
	FoxAddr    string
	FoxDialTls bool
}

func CheckDocumentResError(res *protos.DocumentRes, err error) error {
	if err != nil {
		return fmt.Errorf("gRPC Internal Error: %s", err.Error())
	}
	if res == nil {
		return fmt.Errorf("UnexpectedResponse: gRPC res is null")
	}
	return CheckCommonResError(res.CommonRes, err)
}

func CheckCommonResError(res *protos.CommonRes, err error) error {
	if err != nil {
		return fmt.Errorf("gRPC Internal Error: %s", err.Error())
	}
	if res == nil {
		return fmt.Errorf("UnexpectedResponse: gRPC CommonRes is null")
	}
	if res.Status != protos.ResultCode_SUCCESS {
		return fmt.Errorf("UnexpectedResponse: %s", res.Message)
	}
	return nil
}

func CheckRpcNotExists(res *protos.DocumentRes) bool {
	return res != nil && res.CommonRes != nil && res.CommonRes.Status == protos.ResultCode_NOT_FOUND
}
