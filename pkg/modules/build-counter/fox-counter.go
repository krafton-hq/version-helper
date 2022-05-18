package build_counter

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.krafton.com/xtrm/fox/client/fox_grpc"
	"github.krafton.com/xtrm/fox/core/generated/protos"
)

const MaxUint = ^uint(0)

type FoxCounter struct {
	name              string
	foxServiceProject string
	foxClient         *fox_grpc.FoxClient

	cachedCount uint
}

func NewFoxCounter(name string, foxServiceProject string, foxClient *fox_grpc.FoxClient) *FoxCounter {
	return &FoxCounter{
		name:              name,
		foxServiceProject: foxServiceProject,
		foxClient:         foxClient,
		cachedCount:       MaxUint,
	}
}

func (c *FoxCounter) String() string {
	return fmt.Sprintf("%#v", c)
}

func (c *FoxCounter) Increase(ctx context.Context) (uint, error) {
	res, err := c.foxClient.PatchDocumentField(ctx, &protos.PatchDocumentFieldReq{
		Id:             c.name,
		ServiceProject: c.foxServiceProject,
		Expression:     ".count = .count + 1",
	})
	err = checkRpcError(res, err)
	if err != nil {
		return 0, fmt.Errorf("IncraseFailed: %s", err.Error())
	}

	count, err := c.getCountFromDocument(res.Document)
	if err != nil {
		return 0, fmt.Errorf("IncraseFailed : %s", err.Error())
	}

	c.cachedCount = count
	return count, nil
}

func (c *FoxCounter) Get(ctx context.Context) (uint, error) {
	if c.cachedCount != MaxUint {
		return c.cachedCount, nil
	}

	// Fetch Count From Server
	res, err := c.foxClient.GetDocument(ctx, &protos.GetDocumentReq{
		Id:             c.name,
		ServiceProject: c.foxServiceProject,
	})
	err = checkRpcError(res, err)
	if err != nil {
		return 0, fmt.Errorf("GetFailed: %s", err.Error())
	}

	count, err := c.getCountFromDocument(res.Document)
	if err != nil {
		return 0, fmt.Errorf("GetFailed : %s", err.Error())
	}

	c.cachedCount = count
	return count, nil
}

func (c *FoxCounter) getCountFromDocument(document *protos.DetailedDocument) (uint, error) {
	if document == nil {
		return 0, fmt.Errorf("InvalidParameters: 'document' should not be null")
	}

	rawStruct := map[string]interface{}{}
	err := json.Unmarshal([]byte(document.Document.RawData), &rawStruct)
	if err != nil {
		return 0, fmt.Errorf("JsonParseError: %s", err.Error())
	}

	if value, ok := rawStruct["count"]; ok {
		count, err := strconv.ParseUint(value.(string), 10, strconv.IntSize)
		if err != nil {
			return 0, err
		}
		return uint(count), nil
	} else {
		return 0, fmt.Errorf("InvalidDataFormat: 'count' field does not exists in data, '%s'", document.Document.RawData)
	}
}

func checkRpcError(res *protos.DocumentRes, err error) error {
	if err != nil {
		return fmt.Errorf("gRPC Internal Error: %s", err.Error())
	}
	if res.CommonRes == nil {
		return fmt.Errorf("UnexpectedResponse: gRPC CommonRes is null")
	}
	if res.CommonRes.Status != protos.ResultCode_SUCCESS {
		return fmt.Errorf("UnexpectedResponse: %s", res.CommonRes.Message)
	}
	return nil
}
