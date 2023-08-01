package server

import (
	"net/http"
)

type (
	Handler http.Handler

	RouteList map[string]Handler
)

func NewHandler(routes RouteList) (Handler, error) {
	var mux = http.NewServeMux()
	for route, handler := range routes {
		mux.Handle(route, handler)
	}

	return mux, nil
}
