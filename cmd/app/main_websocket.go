package main

import (
	webSocketHandler "core-finance-ledger/internal/adapters/api/websocket"
	"core-finance-ledger/internal/adapters/cache"
	"log"
	"net/http"
)

type WebSocketMain struct {
	contextPath string
	port        string
}

func NewWebSocketMain(contextPath string, port string) *WebSocketMain {
	return &WebSocketMain{
		contextPath: contextPath,
		port:        port,
	}
}

func (ws *WebSocketMain) StartWebSocketServer(currencies chan map[string]float64, redis *cache.RedisCache) error {
	wsHandler := webSocketHandler.NewWebSocketHandler(currencies, redis)
	go wsHandler.Broadcast()

	http.HandleFunc(ws.contextPath+"/purchase-options/broadcast", wsHandler.HandleConnections)

	log.Println("WebSocket server started on " + ws.port)
	return http.ListenAndServe(ws.port, nil)
}
