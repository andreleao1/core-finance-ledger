package handlers

import (
	"net/http"
	"pricing-simulator/domain/server"
)

func WebSocketHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srv.HandleConnections(w, r)
	}
}
