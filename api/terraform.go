package api

import (
	"net/http"

	"github.com/jaehong21/icarus-proxy/github"
)

func CreateCloudflareTerraformResource(w http.ResponseWriter, r *http.Request) {
	var body NameDto
	if err := ParseJSON(r, &body); err != nil {
		JSON(w, err, http.StatusBadRequest)
		return
	}

	err := github.CreateRoute53(body.Name)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, "success", http.StatusCreated)
}

func DeleteCloudflareTerraformResource(w http.ResponseWriter, r *http.Request) {
	var body NameDto
	if err := ParseJSON(r, &body); err != nil {
		JSON(w, err, http.StatusBadRequest)
		return
	}

	err := github.DeleteRoute53(body.Name)
	if err != nil {
		JSON(w, err, http.StatusInternalServerError)
		return
	}

	JSON(w, "success", http.StatusCreated)
}
