package build_counter

import (
	"context"
	"fmt"

	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	redfoxClientset "github.com/krafton-hq/redfox/pkg/generated/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type RedFoxCounter struct {
	redfoxClient redfoxClientset.Interface
	namespace    string
	name         string

	cachedCount uint
}

func NewRedFoxCounter(redfoxClient redfoxClientset.Interface, namespace string, name string) *RedFoxCounter {
	return &RedFoxCounter{
		redfoxClient: redfoxClient,
		namespace:    namespace,
		name:         name,

		cachedCount: MaxUint,
	}
}

func (c *RedFoxCounter) String() string {
	return fmt.Sprintf("%#v", c)
}

var versioncountKind = schema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "VersionCount"}

const versionHelperManager = "version-helper-cli"

func (c *RedFoxCounter) Increase(ctx context.Context) (uint, error) {
	client := c.redfoxClient.MetadataV1alpha1().VersionCounts(c.namespace)

	count, err := client.Get(ctx, c.name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return 0, fmt.Errorf("GetPreviousCount from Server Failed: %s", err.Error())
	}

	if errors.IsNotFound(err) {
		count = &redfoxV1alpha1.VersionCount{
			TypeMeta: metav1.TypeMeta{
				Kind:       versioncountKind.Kind,
				APIVersion: versioncountKind.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: c.name,
			},
			Spec: redfoxV1alpha1.VersionCountSpec{
				ProjectName: c.name,
				Count:       0,
			},
		}
		count.Spec.Count++

		_, err = client.Create(ctx, count, metav1.CreateOptions{FieldManager: versionHelperManager})
		if err != nil {
			return 0, fmt.Errorf("CreateCount to Server Failed: %s", err.Error())
		}
	} else {
		count.TypeMeta = metav1.TypeMeta{
			Kind:       versioncountKind.Kind,
			APIVersion: versioncountKind.GroupVersion().String(),
		}
		count.Spec.Count++

		_, err = client.Update(ctx, count, metav1.UpdateOptions{FieldManager: versionHelperManager})
		if err != nil {
			return 0, fmt.Errorf("UpdateCount to Server Failed: %s", err.Error())
		}
	}

	c.cachedCount = uint(count.Spec.Count)
	return uint(count.Spec.Count), nil
}

func (c *RedFoxCounter) Get(ctx context.Context) (uint, error) {
	if c.cachedCount != MaxUint {
		return c.cachedCount, nil
	}

	client := c.redfoxClient.MetadataV1alpha1().VersionCounts(c.namespace)

	count, err := client.Get(ctx, c.name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("GetPreviousCount from Server Failed: %s", err.Error())
		}
	}

	return uint(count.Spec.Count), nil
}
