package fox_utils

import (
	"os"
	"path/filepath"

	"github.com/krafton-hq/red-fox/pkg/generated/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const versionHelperUserAgent = "version-helper-cli"

func NewRedFoxClient() (versioned.Interface, error) {
	home, _ := os.UserHomeDir()
	localKubeconfigPath := filepath.Join(home, ".kube", "config")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", localKubeconfigPath)
	if err != nil {
		zap.S().Debugw("Read Kubeconfig failed", "path", localKubeconfigPath, "error", err)
		return nil, err
	}
	zap.S().Debugf("Select %s k8s apiserver", kubeconfig.Host)

	redfoxClient, err := versioned.NewForConfig(rest.AddUserAgent(kubeconfig, versionHelperUserAgent))
	if err != nil {
		zap.S().Debugw("Create redfox client failed", "error", err)
		return nil, err
	}
	return redfoxClient, nil
}
