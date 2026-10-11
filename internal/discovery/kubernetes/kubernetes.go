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

func getPodImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting pods: %w", err)
	}
	for _, pod := range pods.Items {
		for _, i := range pod.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range pod.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range pod.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getDeploymentImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	deployments, err := client.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting deployments: %w", err)
	}
	for _, deployment := range deployments.Items {
		for _, i := range deployment.Spec.Template.Spec.InitContainers {
			res = append(res, i.Image)
		}
		for _, i := range deployment.Spec.Template.Spec.Containers {
			res = append(res, i.Image)
		}
		for _, i := range deployment.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Image)
		}
	}
	return res, nil
}

func getStatefulsetImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	statefulsets, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting statefulsets: %w", err)
	}
	for _, statefulset := range statefulsets.Items {
		for _, i := range statefulset.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range statefulset.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range statefulset.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getDaemonsetImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	daemonsets, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting daemonsets: %w", err)
	}
	for _, daemonset := range daemonsets.Items {
		for _, i := range daemonset.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range daemonset.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range daemonset.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getCronjobImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	cronjobs, err := client.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting cronjobs: %w", err)
	}
	for _, cronjob := range cronjobs.Items {
		for _, i := range cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range cronjob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range cronjob.Spec.JobTemplate.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getJobImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	jobs, err := client.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting jobs: %w", err)
	}
	for _, job := range jobs.Items {
		for _, i := range job.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range job.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range job.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getReplicasetImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	replicasets, err := client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting replicasets: %w", err)
	}
	for _, replicaset := range replicasets.Items {
		for _, i := range replicaset.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range replicaset.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range replicaset.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}

func getReplicationControllerImages(client *kubernetes.Clientset, namespace string) ([]string, error) {
	res := []string{}
	replicationControllers, err := client.CoreV1().ReplicationControllers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return res, fmt.Errorf("error in getting replication controllers: %w", err)
	}
	for _, replicationController := range replicationControllers.Items {
		for _, i := range replicationController.Spec.Template.Spec.Containers {
			res = append(res, i.Name)
		}
		for _, i := range replicationController.Spec.Template.Spec.InitContainers {
			res = append(res, i.Name)
		}
		for _, i := range replicationController.Spec.Template.Spec.EphemeralContainers {
			res = append(res, i.Name)
		}
	}
	return res, nil
}
