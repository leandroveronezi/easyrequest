package easyRequest

import (
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

func init() {

	bindings = make(map[string]bindingType)
	wsBindings = make(map[string]wsBindingType)

	http.HandleFunc("/easyRequest", HandleEasyRequest)
	http.HandleFunc("/wsEasyRequest", WSHandleEasyRequest)
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

func WSHandleEasyRequest(w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, jsString("Url Param 'name' is missing"))
		return
	}

	wsbinding, ok := wsBindings[name[0]]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, jsString("Not Found"))
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, jsString(err.Error()))
		return
	}

	req := WsRequest{}
	req.R = r
	req.W = w
	req.Conn = conn

	result := ""

	err = wsbinding.wsBindingFunc(req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result = jsString(err.Error())
	} else {
		result = string("")
		w.WriteHeader(http.StatusOK)
	}

	req.Conn.Close()

	io.WriteString(w, result)

}
