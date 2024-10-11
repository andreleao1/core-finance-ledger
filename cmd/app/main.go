package main

import (
	webSocketHandler "core-finance-ledger/internal/adapters/api/websocket"
	"core-finance-ledger/internal/adapters/cache"
	"core-finance-ledger/internal/domain/usecase"
	"log"
	"net/http"
)

func main() {
	currencies := make(chan map[string]float64)

	// Inicia o simulador de preços
	redis := cache.NewRedisCache()
	priceSimulator := usecase.NewBitcoinUsecase(redis)
	go priceSimulator.StartPriceSimulation(currencies)

	// Inicia o WebSocket handler
	wsHandler := webSocketHandler.NewWebSocketHandler(currencies, redis)
	go wsHandler.Broadcast()

	http.HandleFunc("/ws", wsHandler.HandleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
