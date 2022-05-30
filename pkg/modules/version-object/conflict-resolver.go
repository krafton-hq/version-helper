package version_object

import (
	"context"
	"fmt"

	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	"github.krafton.com/xtrm/fox/client/fox_grpc"
	"github.krafton.com/xtrm/fox/core/generated/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/yaml"
)

type ConflictResolver interface {
	String() string
	Resolve(ctx context.Context, obj *VersionObj) (*VersionObj, *timestamppb.Timestamp, error)
}

type MergeResolver struct {
	client    *fox_grpc.FoxClient
	namespace string
}

func NewMergeResolver(client *fox_grpc.FoxClient, namespace string) *MergeResolver {
	return &MergeResolver{
		client:    client,
		namespace: namespace,
	}
}

func (r *MergeResolver) String() string {
	return "MergeResolver"
}

func (r *MergeResolver) Resolve(ctx context.Context, obj *VersionObj) (*VersionObj, *timestamppb.Timestamp, error) {
	res, err := r.client.GetDocument(ctx, &protos.GetDocumentReq{
		Id:             obj.Metadata.Name,
		ServiceProject: r.namespace,
	})
	if err := fox_utils.CheckDocumentResError(res, err); err != nil {
		return nil, nil, fmt.Errorf("GetFailed: %s", err.Error())
	}

	doc := res.Document.Document

	oldObj := &VersionObj{}
	err = yaml.Unmarshal([]byte(doc.RawData), oldObj)
	if err != nil {
		return nil, nil, err
	}
	newObj, err := Merge(oldObj, obj)
	if err != nil {
		return nil, nil, err
	}

	return newObj, res.Document.LastModified, nil
}

type OverwriteResolver struct {
}

func NewOverwriteResolver() *OverwriteResolver {
	return &OverwriteResolver{}
}

func (r *OverwriteResolver) String() string {
	return "OverwriteResolver"
}

func (r *OverwriteResolver) Resolve(ctx context.Context, obj *VersionObj) (*VersionObj, *timestamppb.Timestamp, error) {
	return obj, nil, nil
}
