package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/router"
)

func New(app *gtea.App) http.Handler {
	r := router.DefaultRoutesGroup()
	addBusinessRoutes(app, r)
	return r
}

func addBusinessRoutes(app *gtea.App, r *router.RoutesGroup) {
	host := &hostHandler{app: app}
	g := r.Group("/api/v1")
	g.Get("/", hw)
	a := g.Group("/auth", handler.TokenMiddleware(auth))
	a.Get("/hosts", host.list)

}
