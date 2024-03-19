package web

import (
	"net/http"
)

// WebHandler is an interface for handling web requests.
type WebHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
