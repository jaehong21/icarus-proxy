package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/k8s"
	"k8s.io/client-go/kubernetes"
)

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
