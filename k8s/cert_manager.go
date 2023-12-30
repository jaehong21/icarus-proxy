package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func GetCertificates(client dynamic.Interface, namespace string) (*unstructured.UnstructuredList, error) {
	certGVR := schema.GroupVersionResource{
		Group:    "cert-manager.io",
		Version:  "v1",
		Resource: "certificates",
	}

	certs, err := client.Resource(certGVR).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// for _, cert := range certs.Items {
	// 	fmt.Printf("Name: %s, Namespace: %s\n", cert.GetName(), cert.GetNamespace())
	// }

	return certs, err
}
