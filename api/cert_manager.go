package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/k8s"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func GetCertificates(w http.ResponseWriter, r *http.Request, dynamicClient *dynamic.DynamicClient, client *kubernetes.Clientset) {
	namespace := mux.Vars(r)["namespace"]

	nsList, err := k8s.GetNamespaces(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}
	if res := namespaceExists(nsList, namespace); !res {
		JSON(w, "namespace `"+namespace+"` not exists", http.StatusBadRequest)
		return
	}

	certs, err := k8s.GetCertificates(dynamicClient, namespace)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, certs, http.StatusOK)
}
