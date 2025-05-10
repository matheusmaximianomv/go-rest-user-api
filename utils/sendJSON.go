package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SendJSON(w http.ResponseWriter, response Response, statusCode int) error {
	data, err := json.Marshal(response)
	if err != nil {
		SendJSON(w, Response{Message: "Does not convert data"}, http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(statusCode)
	if _, err = w.Write(data); err != nil {
		return fmt.Errorf("does not send response: %w", err)
	}

	return nil
}
