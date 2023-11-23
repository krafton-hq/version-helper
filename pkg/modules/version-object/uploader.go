package version_object

import (
	"context"
	"encoding/json"
	"fmt"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	redfoxClientset "github.com/krafton-hq/redfox/pkg/generated/clientset/versioned"
	"github.com/krafton-hq/version-helper/pkg/modules/versions"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type Uploader struct {
	redfoxClient redfoxClientset.Interface
	namespace    string
}

func NewUploader(redfoxClient redfoxClientset.Interface, namespace string) *Uploader {
	return &Uploader{redfoxClient: redfoxClient, namespace: namespace}
}

const versionHelperManager = "version-helper-cli"

var latestversionsKind = schema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "LatestVersion"}
var versionKind = schema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "LatestVersion"}

func (u *Uploader) Upload(ctx context.Context, version *redfoxV1alpha1.Version) error {
	latestVersion := u.generateLatestVersion(version)

	latestVersion, err := u.ApplyLatestVersion(ctx, latestVersion)
	if err != nil {
		return err
	}

	apiVersion, kind := latestversionsKind.ToAPIVersionAndKind()
	blockOwnerDeletion := true
	isController := true
	version.SetOwnerReferences([]metav1.OwnerReference{
		{
			APIVersion:         apiVersion,
			Kind:               kind,
			Name:               latestVersion.Name,
			UID:                latestVersion.UID,
			Controller:         &isController,
			BlockOwnerDeletion: &blockOwnerDeletion,
		},
	})
	_, err = u.applyVersion(ctx, version)
	return err
}

func (u *Uploader) generateLatestVersion(version *redfoxV1alpha1.Version) *redfoxV1alpha1.LatestVersion {
	var latestVersionName = ""

	if version.Spec.VersionDetail.SubProjectName == "" {
		latestVersionName = fmt.Sprintf("%s-%s", version.Spec.VersionDetail.ProjectName, versions.MangleBranch(version.Spec.GitRef.Branch))
	} else {
		latestVersionName = fmt.Sprintf("%s-%s-%s", version.Spec.VersionDetail.ProjectName, version.Spec.VersionDetail.SubProjectName, versions.MangleBranch(version.Spec.GitRef.Branch))
	}

	latestVersion := &redfoxV1alpha1.LatestVersion{
		TypeMeta: metav1.TypeMeta{
			Kind:       latestversionsKind.Kind,
			APIVersion: version.TypeMeta.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: latestVersionName,
			Labels: map[string]string{
				"repository": version.Spec.GitRef.Repository,
				"branch":     versions.MangleBranch(version.Spec.GitRef.Branch),
			},
		},
		Spec: redfoxV1alpha1.LatestVersionSpec{
			GitRef: redfoxV1alpha1.LatestVersionGitRef{
				Branch:     version.Spec.GitRef.Branch,
				Repository: version.Spec.GitRef.Repository,
			},
		},
		Status: redfoxV1alpha1.LatestVersionStatus{
			VersionRef: redfoxV1alpha1.LatestVersionVersionRef{
				Name: version.Name,
			},
		},
	}
	return latestVersion
}

func (u *Uploader) applyVersion(ctx context.Context, version *redfoxV1alpha1.Version) (*redfoxV1alpha1.Version, error) {
	buf, err := json.Marshal(version)
	if err != nil {
		zap.S().Debugw("Failed to marshal Version object", "error", err, "name", version.Name)
		return nil, err
	}

	_, err = u.redfoxClient.MetadataV1alpha1().Versions(u.namespace).Patch(ctx, version.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: versionHelperManager})
	if err != nil {
		zap.S().Debugw("Failed to apply Version object", "error", err, "name", version.Name)
		return nil, err
	}

	obj, err := u.redfoxClient.MetadataV1alpha1().Versions(u.namespace).Patch(ctx, version.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: versionHelperManager}, "status")
	if err != nil {
		zap.S().Debugw("Failed to apply status Version object", "error", err, "name", version.Name)
		return nil, err
	}

	zap.S().Debugw("Version Object Apply Success", "applied-object", obj)
	return obj, nil
}

func (u *Uploader) ApplyLatestVersion(ctx context.Context, latestVersion *redfoxV1alpha1.LatestVersion) (*redfoxV1alpha1.LatestVersion, error) {
	buf, err := json.Marshal(latestVersion)
	if err != nil {
		zap.S().Debugw("Failed to marshal LatestVersion object", "error", err, "name", latestVersion.Name)
		return nil, err
	}

	forcePatchOption := true
	_, err = u.redfoxClient.MetadataV1alpha1().LatestVersions(u.namespace).Patch(ctx, latestVersion.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: versionHelperManager, Force: &forcePatchOption})
	if err != nil {
		zap.S().Debugw("Failed to apply status LatestVersion object", "error", err, "name", latestVersion.Name)
		return nil, err
	}

	obj, err := u.redfoxClient.MetadataV1alpha1().LatestVersions(u.namespace).Patch(ctx, latestVersion.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: versionHelperManager}, "status")
	if err != nil {
		zap.S().Debugw("Failed to apply LatestVersion object", "error", err, "name", latestVersion.Name)
		return nil, err
	}

	zap.S().Debugw("LatestVersion Object Apply Success", "applied-object", obj)
	return obj, nil
}
