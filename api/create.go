package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/k8s"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset) {
	namespaceName := mux.Vars(r)["namespace"]

	nsList, err := k8s.GetNamespaces(client)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}
	if res := namespaceExists(nsList, namespaceName); res {
		JSON(w, "Namespace `"+namespaceName+"` Already Exists", http.StatusBadRequest)
		return
	}

	ns, err := k8s.CreateNamespace(client, namespaceName)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, ns, http.StatusCreated)
}
