package version_object_service

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	redfoxScheme "github.com/krafton-hq/redfox/pkg/generated/clientset/versioned/scheme"
	fox_utils "github.com/krafton-hq/version-helper/pkg/modules/fox-utils"
	path_utils "github.com/krafton-hq/version-helper/pkg/modules/path-utils"
	version_object "github.com/krafton-hq/version-helper/pkg/modules/version-object"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var redfoxDecoder = redfoxScheme.Codecs.UniversalDeserializer()

func LoadVersionObj(path string) (*redfoxV1alpha1.Version, error) {
	absPath, err := path_utils.ResolvePathToAbs(path)
	if err != nil {
		zap.S().Debugw("Resolve Absolute Path Failed", "path", absPath, "error", err)
		return nil, err
	}
	buf, err := ioutil.ReadFile(absPath)
	if err != nil {
		zap.S().Debugw("Read Version Object Failed", "path", absPath, "error", err)
		return nil, err
	}

	version := &redfoxV1alpha1.Version{}
	if _, _, err = redfoxDecoder.Decode(buf, nil, version); err != nil {
		zap.S().Debugw("Marshal raw data to version object failed", "error", err)
	}
	return version, nil
}

func SaveVersionObj(obj *redfoxV1alpha1.Version, path string) error {
	buffer := &bytes.Buffer{}
	serializer := json.NewSerializerWithOptions(yaml.DefaultMetaFactory, redfoxScheme.Scheme, redfoxScheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: true, Strict: false})
	err := serializer.Encode(obj, buffer)
	if err != nil {
		zap.S().Debugw("Unmarshal version object to raw data failed", "error", err)
		return err
	}

	absPath, err := path_utils.ResolvePathToAbs(path)
	if err != nil {
		zap.S().Debugw("Resolve Absolute Path Failed", "path", absPath, "error", err)
		return err
	}

	err = ioutil.WriteFile(absPath, buffer.Bytes(), 0644)
	if err != nil {
		zap.S().Infof("Write File Failed, error: %s", err.Error())
		return err
	}
	return nil
}

const redfoxNamespace = "redfox-metadata"

func UploadVersionObj(ctx context.Context, obj *redfoxV1alpha1.Version) error {
	redFoxClient, err := fox_utils.NewRedFoxClient()
	if err != nil {
		return err
	}

	uploader := version_object.NewUploader(redFoxClient, redfoxNamespace)
	return uploader.Upload(ctx, obj)
}

func UploadLatestVersion(ctx context.Context, obj *redfoxV1alpha1.LatestVersion, namespace string) (*redfoxV1alpha1.LatestVersion, error) {
	redFoxClient, err := fox_utils.NewRedFoxClient()
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		return nil, errors.New("namespace is empty")
	}

	uploader := version_object.NewUploader(redFoxClient, namespace)
	return uploader.ApplyLatestVersion(ctx, obj)
}

func GetLatestVersion(ctx context.Context, name, namespace string) (*redfoxV1alpha1.LatestVersion, error) {
	redFoxClient, err := fox_utils.NewRedFoxClient()
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		return nil, errors.New("namespace is empty")
	}

	lv, err := redFoxClient.MetadataV1alpha1().LatestVersions(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return lv, nil
}
