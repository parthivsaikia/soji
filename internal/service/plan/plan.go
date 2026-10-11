package plan

import "github.com/parthivsaikia/soji/internal/discovery/kubernetes"

func Plan() ([]string, error) {
	kubeconfig := "/home/parthiv/.kube/config"
	images, err := kubernetes.FindCurrentlyUsedImage(kubeconfig)
	if err != nil {
		return images, err
	}
	return images, nil
}
