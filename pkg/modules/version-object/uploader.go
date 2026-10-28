package version_object

import (
	"context"
	"encoding/json"
	"fmt"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	redfoxClientset "github.com/krafton-hq/redfox/pkg/generated/clientset/versioned"
	"github.com/krafton-hq/version-helper/pkg/modules/versions"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type Uploader struct {
	redfoxClient redfoxClientset.Interface
	namespace    string
	// forceLatest 가 true 면 기존 LatestVersion 이 더 높은 revision 을 가리켜도 이번 버전으로 덮어쓴다(예전 동작).
	forceLatest bool
}

func NewUploader(redfoxClient redfoxClientset.Interface, namespace string) *Uploader {
	return &Uploader{redfoxClient: redfoxClient, namespace: namespace}
}

// WithForceLatest 는 LatestVersion 을 revision 비교 없이 항상 덮어쓰게 한다.
func (u *Uploader) WithForceLatest(force bool) *Uploader {
	u.forceLatest = force
	return u
}

const versionHelperManager = "version-helper-cli"

var latestversionsKind = schema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "LatestVersion"}
var versionKind = schema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "LatestVersion"}

func (u *Uploader) Upload(ctx context.Context, version *redfoxV1alpha1.Version) error {
	latestVersion := u.generateLatestVersion(version)

	current, keep, err := u.keepCurrentLatest(ctx, latestVersion.Name, version)
	if err != nil {
		return err
	}
	if keep {
		// 같은 브랜치의 더 높은 revision 이 이미 latest 다. 지난 버전을 다시 올리거나 순서가 뒤바뀐 재시도가
		// latest 포인터를 과거로 되돌리지 않게, Version 만 올리고 LatestVersion 은 그대로 둔다.
		zap.S().Infow("LatestVersion already points to a newer revision; keeping it",
			"latestVersion", current.Name, "current", current.Status.VersionRef.Name, "uploading", version.Name)
		latestVersion = current
	} else {
		latestVersion, err = u.ApplyLatestVersion(ctx, latestVersion)
		if err != nil {
			return err
		}
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

// keepCurrentLatest 는 기존 LatestVersion 을 유지해야 하는지 판단한다.
// 기존 LatestVersion 이 가리키는 Version 의 spec.versionDetail.revision 이 올리는 버전보다 크면 유지한다.
// revision 이 없는(0) 버전끼리는 비교할 근거가 없으므로 예전처럼 덮어쓴다.
func (u *Uploader) keepCurrentLatest(ctx context.Context, latestVersionName string, version *redfoxV1alpha1.Version) (*redfoxV1alpha1.LatestVersion, bool, error) {
	if u.forceLatest {
		return nil, false, nil
	}
	current, err := u.redfoxClient.MetadataV1alpha1().LatestVersions(u.namespace).Get(ctx, latestVersionName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		zap.S().Debugw("Failed to get LatestVersion object", "error", err, "name", latestVersionName)
		return nil, false, err
	}
	refName := current.Status.VersionRef.Name
	if refName == "" || refName == version.Name {
		return current, false, nil
	}
	ref, err := u.redfoxClient.MetadataV1alpha1().Versions(u.namespace).Get(ctx, refName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return current, false, nil
	}
	if err != nil {
		zap.S().Debugw("Failed to get referenced Version object", "error", err, "name", refName)
		return nil, false, err
	}
	return current, isNewerRevision(ref.Spec.VersionDetail.Revision, version.Spec.VersionDetail.Revision), nil
}

// isNewerRevision 은 현재 latest 의 revision 이 올리는 버전보다 큰지 본다.
// 어느 한쪽이라도 revision 이 없으면(0) 비교 근거가 없으므로 false(=덮어쓰기)다.
func isNewerRevision(currentRevision, uploadingRevision int) bool {
	if currentRevision == 0 || uploadingRevision == 0 {
		return false
	}
	return currentRevision > uploadingRevision
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
