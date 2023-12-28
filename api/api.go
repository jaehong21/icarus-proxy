package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/k8s"
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

func GetPods(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespace := mux.Vars(r)["namespace"]

	pods, err := k8s.GetPods(client, namespace)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, pods, http.StatusOK)
}