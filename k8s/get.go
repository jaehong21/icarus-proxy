package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNodes(client *kubernetes.Clientset) (*v1.NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetNamespaces(client *kubernetes.Clientset) (*v1.NamespaceList, error) {
	ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// GetNamespaceStatus gets the status of a specific namespace
func GetNamespaceStatus(client *kubernetes.Clientset, namespaceName string) (string, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Check if the namespace is being deleted
	if ns.Status.Phase == v1.NamespaceTerminating {
		return "Terminating", nil
	}

	// If the namespace exists and is not in terminating state, it's considered 'Active'
	// if ns.Status.Phase == v1.NamespaceActive
	return "Active", nil
}

func GetPods(client *kubernetes.Clientset, namespace string) (*v1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}
