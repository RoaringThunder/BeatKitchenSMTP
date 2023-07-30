package utils

import (
	"encoding/json"
	"net/http"
	"salamander-smtp/logging"
)

func HTTPHandleResponse(w http.ResponseWriter, payload map[string]interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		logging.Log("Failed marshalling json response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(200)
	w.Write([]byte(response))
}

func HTTPHandleError(w http.ResponseWriter, code int, errorStr string) {
	payload := map[string]interface{}{
		"status":  false,
		"message": errorStr,
	}
	response, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		logging.Log("Failed to marshall response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("We're having some issues right now"))
		return
	}
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}
