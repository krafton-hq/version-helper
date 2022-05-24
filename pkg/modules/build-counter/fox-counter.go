package build_counter

import (
	"context"
	"encoding/json"
	"fmt"

	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
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
	err = fox_utils.CheckRpcError(res, err)
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
	if fox_utils.CheckRpcNotExists(res) {
		res, err := c.foxClient.CreateDocument(ctx, &protos.CreateDocumentReq{
			Document: &protos.Document{
				Id:         c.name,
				RawData:    "{\"count\": 1}",
				DataType:   "map",
				ApiVersion: "version.sbx-central.io/v1alpha1",
				Kind:       "Counter",
			},
			ServiceProject: c.foxServiceProject,
		})
		if err := fox_utils.CheckCommonRpcError(res, err); err != nil {
			return 0, fmt.Errorf("GetFailed: %s", err.Error())
		}

		c.cachedCount = 1
		return c.cachedCount, nil
	}

	if err := fox_utils.CheckRpcError(res, err); err != nil {
		return 0, fmt.Errorf("GetFailed: %s", err.Error())
	}

	count, err := c.getCountFromDocument(res.Document)
	if err != nil {
		return 0, fmt.Errorf("GetFailed : %s", err.Error())
	}

	c.cachedCount = count
	return count, nil
}

type countData struct {
	Count uint `json:"count"`
}

func (c *FoxCounter) getCountFromDocument(document *protos.DetailedDocument) (uint, error) {
	if document == nil {
		return 0, fmt.Errorf("InvalidParameters: 'document' should not be null")
	}

	countStruct := &countData{}
	err := json.Unmarshal([]byte(document.Document.RawData), &countStruct)
	if err != nil {
		return 0, fmt.Errorf("JsonParseError: %s", err.Error())
	}

	return countStruct.Count, nil
}
