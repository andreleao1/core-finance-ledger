package main

import (
	webSocketHandler "core-finance-ledger/internal/adapters/api/websocket"
	"core-finance-ledger/internal/adapters/cache"
	"core-finance-ledger/internal/domain/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const CONTEXT_PATH = "/core-finance-ledger"
const WEBSOCKET_SERVER_PORT = ":9096"
const HTTP_SERVER_PORT = ":9097"

func main() {
	currencies := make(chan map[string]float64)

	// Inicia o simulador de preços
	redis := cache.NewRedisCache()
	priceSimulator := usecase.NewBitcoinUsecase(redis)
	go priceSimulator.StartPriceSimulation(currencies)

	startHttpServer()

	// Inicia o WebSocket handler
	startWebSocketServer(currencies, redis)
}

func startWebSocketServer(currencies chan map[string]float64, redis *cache.RedisCache) {
	wsHandler := webSocketHandler.NewWebSocketHandler(currencies, redis)
	go wsHandler.Broadcast()

	http.HandleFunc(CONTEXT_PATH+"/purchase-options/broadcast", wsHandler.HandleConnections)

	log.Println("Websocket server started on " + WEBSOCKET_SERVER_PORT)
	err := http.ListenAndServe(WEBSOCKET_SERVER_PORT, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func startHttpServer() {
	httpServer := gin.Default()
	//contextPath := httpServer.Group(CONTEXT_PATH)

	httpServer.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("HTTP server started on " + HTTP_SERVER_PORT)
	httpServer.Run(HTTP_SERVER_PORT)
}
