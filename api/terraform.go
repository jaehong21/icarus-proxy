package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaehong21/icarus-proxy/github"
)

func CreateCloudflareTerraformResource(w http.ResponseWriter, r *http.Request) {
	dnsName := mux.Vars(r)["name"]
	err := github.CreateRoute53(dnsName)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, "success", http.StatusCreated)
}

func DeleteCloudflareTerraformResource(w http.ResponseWriter, r *http.Request) {
	resourceName := mux.Vars(r)["name"]
	err := github.DeleteRoute53(resourceName)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, "success", http.StatusCreated)
}
