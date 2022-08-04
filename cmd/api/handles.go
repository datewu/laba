package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
)

func hw(w http.ResponseWriter, r *http.Request) {
	msg := "hello world"
	handler.WriteStr(w, http.StatusOK, msg, nil)
}

type hostHandler struct {
	app *gtea.App
}

func (h hostHandler) list(w http.ResponseWriter, r *http.Request) {
	msg := "ping from auth, you've been  authenticated"
	handler.WriteStr(w, http.StatusOK, msg, nil)
}
