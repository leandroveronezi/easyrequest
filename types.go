package easyRequest

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	W http.ResponseWriter
	R *http.Request
}

type bindingFunc func(req Request, args []json.RawMessage) (interface{}, error)

type bindingType struct {
	bindingFunc
	name    string
	comment string
	in      []interface{}
	out     interface{}
}

var bindings map[string]bindingType

type payload struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}
