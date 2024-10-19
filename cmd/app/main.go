package main

import (
	"core-finance-ledger/internal/adapters/cache"
	"core-finance-ledger/internal/domain/usecase"
	"log"
	"sync"
)

func main() {
	currencies := make(chan map[string]float64)

	const CONTEXT_PATH = "/core-finance-ledger"
	const WEBSOCKET_SERVER_PORT = ":9096"
	const HTTP_SERVER_PORT = ":9097"

	// Inicia o simulador de preços
	redis := cache.NewRedisCache()
	priceSimulator := usecase.NewBitcoinUsecase(redis)
	go priceSimulator.StartPriceSimulation(currencies)

	var wg sync.WaitGroup

	httpServer := NewHttpServer(CONTEXT_PATH, HTTP_SERVER_PORT)
	websocketServer := NewWebSocketMain(CONTEXT_PATH, WEBSOCKET_SERVER_PORT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := websocketServer.StartWebSocketServer(currencies, redis); err != nil {
			log.Fatalf("Erro ao iniciar o servidor WebSocket: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.StartHttpServer(); err != nil {
			log.Fatalf("Erro ao iniciar o servidor HTTP: %v", err)
		}
	}()

	// Aguarda a finalização dos servidores
	wg.Wait()
}
