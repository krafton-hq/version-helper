package version_object_service

import (
	"context"
	"fmt"
	"io/ioutil"

	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	path_utils "github.com/krafton-hq/version-helper/pkg/modules/path-utils"
	version_object "github.com/krafton-hq/version-helper/pkg/modules/version-object"
	"github.krafton.com/xtrm/fox/client/fox_grpc"
	"github.krafton.com/xtrm/fox/core/generated/protos"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"
)

func LoadVersionObj(path string) (*version_object.VersionObj, error) {
	absPath, err := path_utils.ResolvePathToAbs(path)
	if err != nil {
		zap.S().Debugf("Resolve Absolute Path Failed, file: %s, error: %s", path, err.Error())
		return nil, err
	}
	buf, err := ioutil.ReadFile(absPath)
	if err != nil {
		zap.S().Debugf("Read Version Object Failed, path: %s, error: %s", absPath, err.Error())
		return nil, err
	}
	obj := &version_object.VersionObj{}
	err = yaml.Unmarshal(buf, obj)
	if err != nil {
		zap.S().Debugw("Deserialize Version Object Failed", "raw", buf, "error", err.Error())
		return nil, err
	}

	return obj, nil
}

func SaveVersionObj(obj *version_object.VersionObj, path string) error {
	buf, err := yaml.Marshal(obj)
	if err != nil {
		zap.S().Debugw("Serialize Version Object Failed", "object", obj, "error", err.Error())
		return err
	}

	absPath, err := path_utils.ResolvePathToAbs(path)
	if err != nil {
		zap.S().Debugf("Resolve Absolute Path Failed, file: %s, error: %s", path, err.Error())
		return err
	}

	err = ioutil.WriteFile(absPath, buf, 0644)
	if err != nil {
		zap.S().Infof("Write File Failed, error: %s", err.Error())
		return err
	}
	return nil
}

// CreateVersionObj returns conflict, error
// Scenarios 1) Normal case (false, nil)
// Scenarios 2) Conflict Detect (true, error)
// Scenarios 3) Other Errors (false, error)
func CreateVersionObj(obj *version_object.VersionObj, client *fox_grpc.FoxClient) (bool, error) {
	return false, nil
}

func ResolveConflictWithRetry(obj *version_object.VersionObj, client *fox_grpc.FoxClient, policy string, retry uint) error {
	var err error
	for i := uint(0); i < retry; i++ {
		err = ResolveConflict(obj, client, policy)
		if err == nil {
			break
		}
	}
	return err
}

const foxNamespace = "SBX-VERSION"

func ResolveConflict(obj *version_object.VersionObj, client *fox_grpc.FoxClient, policy string) error {
	if policy == "merge" {
		res, err := client.GetDocument(context.TODO(), &protos.GetDocumentReq{
			Id:             obj.Metadata.Name,
			ServiceProject: foxNamespace,
		})
		if err := fox_utils.CheckRpcError(res, err); err != nil {
			return fmt.Errorf("GetFailed: %s", err.Error())
		}

		doc := res.Document.Document

		oldObj := &version_object.VersionObj{}
		err = yaml.Unmarshal([]byte(doc.RawData), oldObj)
		if err != nil {
			return err
		}
		newObj, err := version_object.Merge(oldObj, obj)
		if err != nil {
			return err
		}

		buf, err := yaml.Marshal(newObj)
		if err != nil {
			zap.S().Debugw("Serialize Version Object Failed", "object", newObj, "error", err.Error())
			return err
		}

		doc.RawData = string(buf)
		uRes, err := client.UpdateDocument(context.TODO(), &protos.UpdateDocumentReq{
			Document:       doc,
			ServiceProject: foxNamespace,
		})
		if err := fox_utils.CheckCommonRpcError(uRes, err); err != nil {
			return fmt.Errorf("UpdateFailed: %s", err.Error())
		}
		return nil
	} else if policy == "overwrite" {
		//buf, err := yaml.Marshal(obj)
		//if err != nil {
		//	zap.S().Debugw("Serialize Version Object Failed", "object", obj, "error", err.Error())
		//	return err
		//}
		//
		//uRes, err := client.UpdateDocument(context.TODO(), &protos.UpdateDocumentReq{
		//	Document:       doc,
		//	ServiceProject: foxNamespace,
		//})
		//if err := fox_utils.CheckCommonRpcError(uRes, err); err != nil {
		//	return fmt.Errorf("UpdateFailed: %s", err.Error())
		//}
		return nil
	} else {
		return fmt.Errorf("UnexpectedPolicy: %s", policy)
	}
}
