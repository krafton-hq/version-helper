package fox_utils

import (
	"fmt"
	"os"
	"path/filepath"

	redfoxClientset "github.com/krafton-hq/redfox/pkg/generated/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const versionHelperUserAgent = "version-helper-cli"

func NewRedFoxClient() (redfoxClientset.Interface, error) {
	config, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	config = rest.AddUserAgent(config, versionHelperUserAgent)
	return redfoxClientset.NewForConfig(config)
}

func localConfig() (*rest.Config, error) {
	home, _ := os.UserHomeDir()
	localKubeconfigPath := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", localKubeconfigPath)
}

func resolveConfig() (*rest.Config, error) {
	if kubeconfigPath := os.Getenv("REDFOX_KUBECONFIG"); kubeconfigPath != "" {
		// Use Current Context
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	if kubeHost := os.Getenv("REDFOX_HOST"); kubeHost != "" {
		return &rest.Config{Host: kubeHost}, nil
	}

	if config, err := localConfig(); err == nil {
		return config, nil
	}

	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}
	return nil, fmt.Errorf("failed to resolve Kubeconfig for Redfox Client")
}
