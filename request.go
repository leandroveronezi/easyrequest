package easyRequest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func jsString(v interface{}) string {

	b, _ := json.Marshal(v)

	return string(b)

}

func loadRequest(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		//
	}

	var bodyBytes []byte

	bodyBytes, _ = ioutil.ReadAll(r.Body)

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Use the content
	p := payload{}
	err := json.Unmarshal(bodyBytes, &p)

	result := ""

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, jsString(err.Error()))
		return
	}

	binding, ok := bindings[p.Name]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, jsString("Not Found"))

		return
	}

	req := Request{}
	req.R = r
	req.W = w

	if returning, err := binding.bindingFunc(req, p.Args); err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		result = jsString(err.Error())

	} else if b, err := json.Marshal(returning); err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		result = jsString(err.Error())

	} else {

		result = string(b)

		w.WriteHeader(http.StatusOK)

	}

	io.WriteString(w, result)

}
