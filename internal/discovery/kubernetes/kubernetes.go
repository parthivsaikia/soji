package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func FindCurrentlyUsedImage() ([]string, error) {
	kubeconfig := "/home/parthiv/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	images := []string{}
	if err != nil {
		return images, fmt.Errorf("error in creating config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return images, fmt.Errorf("error in creating clientset: %w", err)
	}

	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return images, fmt.Errorf("error in fetching deployments: %w", err)
	}
	for _, deployment := range deployments.Items {
		initContainers := deployment.Spec.Template.Spec.InitContainers
		for _, i := range initContainers {
			images = append(images, i.Image)
		}
		appContainers := deployment.Spec.Template.Spec.Containers
		for _, i := range appContainers {
			images = append(images, i.Image)
		}
		ephemeralContainers := deployment.Spec.Template.Spec.EphemeralContainers
		for _, i := range ephemeralContainers {
			images = append(images, i.Image)
		}
	}
	return images, nil
}
