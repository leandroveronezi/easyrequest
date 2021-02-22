package easyRequest

import (
	"encoding/json"
	"io"
	"net/http"
)

func hdoc(w http.ResponseWriter, r *http.Request) {

	var sliceaux []interface{}

	for _, f := range bindings {

		aux := struct {
			Name    string
			Comment string
			In      []interface{}
			Out     interface{}
		}{
			f.name,
			f.comment,
			f.in,
			f.out,
		}

		sliceaux = append(sliceaux, aux)

	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(sliceaux)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

}
