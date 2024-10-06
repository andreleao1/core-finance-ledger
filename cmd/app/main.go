package main

import (
	webSocketHandler "core-finance-ledger/internal/adapters/api/websocket"
	"core-finance-ledger/internal/domain/usecase"
	"log"
	"net/http"
)

func main() {
	currencies := make(chan map[string]float64)

	// Inicia o simulador de preços
	priceSimulator := usecase.NewPriceSimulator()
	go priceSimulator.StartPriceSimulation(currencies)

	// Inicia o WebSocket handler
	wsHandler := webSocketHandler.NewWebSocketHandler(currencies)
	go wsHandler.Broadcast()

	http.HandleFunc("/ws", wsHandler.HandleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
