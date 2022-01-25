package handlers

import (
	"awesomeProject/models"
	"net/http"
)

func returnJSON(w http.ResponseWriter, obj models.JSONHelper) error {
	w.Header().Add("content-type", "application/json")
	if err := obj.ToJSON(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// returns http error with givven http status code and object
func returnJSONError(writer http.ResponseWriter, obj models.JSONHelper, code int) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(code)
	obj.ToJSON(writer)
}
