package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/k8s"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

func GetNodes(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	nodes, err := k8s.GetNodes(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, nodes, http.StatusOK)
}

func GetNamespaces(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	ns, err := k8s.GetNamespaces(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, ns, http.StatusOK)
}

func GetNamespaceStatus(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespaceName := mux.Vars(r)["namespace"]

	ns, err := k8s.GetNamespaceStatus(client, namespaceName)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, ns, http.StatusOK)
}

func GetPods(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespace := mux.Vars(r)["namespace"]

	pods, err := k8s.GetPods(client, namespace)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, pods, http.StatusOK)
}

func CreateNamespace(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespaceName := mux.Vars(r)["namespace"]

	nsList, err := k8s.GetNamespaces(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}
	if res := namespaceExists(nsList, namespaceName); res {
		JSON(w, "namespace `"+namespaceName+"` already exists", http.StatusBadRequest)
		return
	}

	ns, err := k8s.CreateNamespace(client, namespaceName)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, ns, http.StatusCreated)
}

func DeleteNamespace(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespaceName := mux.Vars(r)["namespace"]

	nsList, err := k8s.GetNamespaces(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}
	if res := namespaceExists(nsList, namespaceName); !res {
		JSON(w, "Namespace `"+namespaceName+"` Not Exists", http.StatusBadRequest)
		return
	}

	if err := k8s.DeleteNamespace(client, namespaceName); err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, "success", http.StatusCreated)
}

func namespaceExists(nsList *v1.NamespaceList, namespaceName string) bool {
	for _, ns := range nsList.Items {
		if ns.Name == namespaceName {
			return true
		}
	}

	return false
}
