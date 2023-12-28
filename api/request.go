package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, value any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")

	w.WriteHeader(statusCode)

	switch v := value.(type) {
	case error:
		resp := map[string]interface{}{
			"msg": v.Error(),
		}
		log.Println("err:", v)
		json.NewEncoder(w).Encode(resp)
	case string:
		resp := map[string]interface{}{
			"msg": v,
		}
		json.NewEncoder(w).Encode(resp)
	default:
		json.NewEncoder(w).Encode(value)
	}
}
