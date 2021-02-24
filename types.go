package easyRequest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type Request struct {
	W http.ResponseWriter
	R *http.Request
}

type WsRequest struct {
	W    http.ResponseWriter
	R    *http.Request
	Conn *websocket.Conn
}

type bindingFunc func(req Request, args []json.RawMessage) (interface{}, error)
type wsBindingFunc func(req WsRequest) error

type bindingType struct {
	bindingFunc
	name string
	in   []interface{}
	out  interface{}
}

type wsBindingType struct {
	wsBindingFunc
	name string
}

var bindings map[string]bindingType
var wsBindings map[string]wsBindingType

type payload struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}
