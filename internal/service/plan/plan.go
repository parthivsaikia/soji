package plan

import "github.com/parthivsaikia/soji/internal/discovery/kubernetes"

func Plan() ([]string, error) {
	images, err := kubernetes.FindCurrentlyUsedImage()
	if err != nil {
		return images, err
	}
	return images, nil
}
