package version_object

import (
	"context"
	"encoding/json"
	"fmt"

	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	"github.krafton.com/xtrm/fox/client/fox_grpc"
	"github.krafton.com/xtrm/fox/core/generated/protos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Uploader struct {
	client           *fox_grpc.FoxClient
	namespace        string
	conflictResolver ConflictResolver
	retry            int
}

func NewUploader(client *fox_grpc.FoxClient, namespace string, conflictResolver ConflictResolver, retry int) *Uploader {
	return &Uploader{client: client, namespace: namespace, conflictResolver: conflictResolver, retry: retry}
}

func (u *Uploader) Upload(ctx context.Context, obj *VersionObj) error {
	// First Try
	err := u.create(ctx, obj)
	if err == nil {
		return nil
	}
	zap.S().Debugw("Document Creation Conflicts to Remote Server", "object", obj, "error", err.Error())

	// Since Second Try
	for i := 1; i < u.retry; i++ {
		zap.S().Debugf("#%d Try to Resolve Conflict Resolver: %s", i, u.conflictResolver.String())
		resolvedObj, cas, err := u.conflictResolver.Resolve(ctx, obj)
		if err != nil {
			zap.S().Debugw("Conflict Resolve Failed", "error", err.Error())
			return err
		}

		err = u.update(ctx, resolvedObj, cas)
		if err == nil {
			zap.S().Debug("Document Update Success")
			return nil
		}
	}

	zap.S().Debugw("Document Creation Failed to Remote Server", "object", obj, "error", err.Error())
	return err
}

func (u *Uploader) create(ctx context.Context, obj *VersionObj) error {
	doc, err := ToFoxDocument(obj)
	if err != nil {
		return err
	}

	res, err := u.client.CreateDocument(ctx, &protos.CreateDocumentReq{
		Document:       doc,
		ServiceProject: u.namespace,
	})
	if err := fox_utils.CheckCommonResError(res, err); err != nil {
		return fmt.Errorf("CreateDocumentFailed: %s", err.Error())
	}
	return nil
}

func (u *Uploader) update(ctx context.Context, obj *VersionObj, cas *timestamppb.Timestamp) error {
	doc, err := ToFoxDocument(obj)
	if err != nil {
		return err
	}

	res, err := u.client.UpdateDocument(ctx, &protos.UpdateDocumentReq{
		Document:        doc,
		ServiceProject:  u.namespace,
		CasLastModified: cas,
	})
	if err := fox_utils.CheckCommonResError(res, err); err != nil {
		return fmt.Errorf("UpdateDocumentFailed: %s", err.Error())
	}
	return nil
}

func ToFoxDocument(obj *VersionObj) (*protos.Document, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		zap.S().Debugw("Serialize Version Object Failed", "object", obj, "error", err.Error())
		return nil, err
	}

	return &protos.Document{
		Id:         obj.Metadata.Name,
		Groups:     LabelsToFoxGroups(obj.Metadata.Labels),
		RawData:    string(buf),
		DataType:   "map",
		ApiVersion: "versions.sbx-central.io/v1alpha1",
		Kind:       "Version",
	}, nil
}

func LabelsToFoxGroups(labels map[string]string) []string {
	var groups []string
	for key, value := range labels {
		groups = append(groups, fmt.Sprintf("%s=", key))
		if value != "" {
			groups = append(groups, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return groups
}
