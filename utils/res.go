package utils

import (
	"encoding/json"
	"net/http"
)

// utilty package to handle response JSON

// p interface{} -> to store request of any type
func ResponseJSON(rw http.ResponseWriter, p interface{}, status int) {
	changeToByte, err := json.Marshal(p)
	rw.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(rw, "Error", http.StatusBadRequest)
	}

	rw.WriteHeader(status)
	rw.Write([]byte(changeToByte))
}
