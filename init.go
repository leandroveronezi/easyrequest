package easyRequest

import (
	"io"
	"net/http"
)

func init() {

	bindings = make(map[string]bindingType)

	http.HandleFunc("/easyRequest", HandleEasyRequest)
}

func HandleEasyRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
		hdoc(w, r)
	case http.MethodPost:
		// Create a new record.
		loadRequest(w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Not found")
		// Give an error message.
	}

}
