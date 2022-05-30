package version_object_service

import (
	"context"
	"fmt"
	"io/ioutil"

	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	path_utils "github.com/krafton-hq/version-helper/pkg/modules/path-utils"
	version_object "github.com/krafton-hq/version-helper/pkg/modules/version-object"
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

const foxNamespace = "SBX-VERSION"

type Option struct {
	// Required
	ConflictResolvePolicy string
	ConflictRetry         int

	// Optional
	FoxAddr    string
	FoxDialTls bool
}

func UploadVersionObj(ctx context.Context, obj *version_object.VersionObj, option *Option) error {
	foxClient, err := fox_utils.NewClient(&fox_utils.Option{
		FoxAddr:    option.FoxAddr,
		FoxDialTls: option.FoxDialTls,
	})
	if err != nil {
		return err
	}

	var conflictResolver version_object.ConflictResolver
	switch option.ConflictResolvePolicy {
	case "merge":
		conflictResolver = version_object.NewMergeResolver(foxClient, foxNamespace)
		break
	case "overwrite":
		conflictResolver = version_object.NewOverwriteResolver()
	default:
		return fmt.Errorf("InvalidArguments, policy: %s", option.ConflictResolvePolicy)
	}

	uploader := version_object.NewUploader(foxClient, foxNamespace, conflictResolver, option.ConflictRetry)
	return uploader.Upload(ctx, obj)
}
