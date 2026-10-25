package plan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parthivsaikia/soji/internal/discovery/kubernetes"
	"github.com/parthivsaikia/soji/internal/discovery/registry"
	"github.com/parthivsaikia/soji/internal/discovery/registry/dockerhub"
	"github.com/parthivsaikia/soji/internal/parser"
)

func Plan(kubeconfig, url string, private bool) error {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error in getting kubeconfig file: %w", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	target, err := parser.Parse(url)
	if err != nil {
		return err
	}
	if target.Provider != parser.ProviderDockerHub {
		return fmt.Errorf("provider not supported")
	}
	dockerhubReg := dockerhub.DockerHubRegistry{
		Namespace:   target.Path,
		ReadPrivate: private,
	}
	registryImages, err := registry.Discover(dockerhubReg)
	fmt.Println("registry images...")
	for _, i := range registryImages {
		fmt.Println(i)
	}
	if err != nil {
		return err
	}
	k8sImages, err := kubernetes.FindCurrentlyUsedImage(kubeconfig)
	fmt.Println("k8s images...")
	for _, i := range k8sImages {
		fmt.Println(i)
	}
	if err != nil {
		return err
	}
	return nil
}
