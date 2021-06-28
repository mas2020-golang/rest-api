package utils

import (
	"fmt"
	"net/http"
)

// ReturnError returns an error to the caller
func ReturnError(w *http.ResponseWriter, message string, responseCode int)  {
	(*w).WriteHeader(responseCode)
	(*w).Header().Set("Content-Type", "application/json")
	body := fmt.Sprintf(`{"error": "%s"}`, message)
	(*w).Write([]byte(body))
}

