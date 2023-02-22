package main

import (
	"context"
	"fmt"

	// k8s.io/api is a set of packages that define Kubernetes API objects
	// - Deployment, Pod, Service, etc.

	// k8s.io/apimachinery is a set of packages that define base-level, universal data structures
	// - ObjectMeta (metadata), TypeMeta (kind, apiVersion), etc.
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Client-go holds the client interfaces
	// - clientcmd sets up a client from a kubeconfig file
	// - kubernetes holds the Kubernetes API clients
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func createDeploymentSpec() appsv1.Deployment {
	deployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "demo", Image: "lisa-agent:latest"},
					},
				},
			},
		},
	}

	fmt.Printf("%#v", deployment)

	return deployment
}

func createJobSpec() batchv1.Job {
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "lisa-job2",
							Image:           "lisa-agent:latest",
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}

	fmt.Printf("%#v", job)
	return job
}

func createGeneralClientSet() *kubernetes.Clientset {
	// uses the current context in kubeconfig to create a client
	// kubeconfig holds client config with server name, credentials, etc
	// Also holds data to configure access to clusters
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\Users\\frmunozp\\.kube\\config")
	if err != nil {
		fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
		panic(err.Error())
	}

	// creates the clientset
	// The client set contains multiple clients for different API groups
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func main() {
	fmt.Println("Deployment:\n=====================\n")
	deploymentSpec := createDeploymentSpec()

	fmt.Println("\nJob:\n=====================\n")
	jobSpec := createJobSpec()

	clientSet := createGeneralClientSet()

	deploymentsClient := clientSet.AppsV1().Deployments(corev1.NamespaceDefault)
	jobsClient := clientSet.BatchV1().Jobs(corev1.NamespaceDefault)

	deploymentsClient.Create(context.TODO(), &deploymentSpec, metav1.CreateOptions{})
	jobsClient.Create(context.TODO(), &jobSpec, metav1.CreateOptions{})
}
