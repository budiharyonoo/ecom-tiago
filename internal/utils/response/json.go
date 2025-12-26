package response

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Json(w http.ResponseWriter, httpStatusCode int, optionalDataAndMessage ...any) {
	var data any = nil
	message := http.StatusText(httpStatusCode)

	if len(optionalDataAndMessage) > 0 {
		data = optionalDataAndMessage[0]
	}
	if len(optionalDataAndMessage) > 1 {
		if msg, ok := optionalDataAndMessage[1].(string); ok {
			message = msg
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := jsonResponse{
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}
