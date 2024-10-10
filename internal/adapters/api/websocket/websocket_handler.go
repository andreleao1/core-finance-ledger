package adapters

import (
	"core-finance-ledger/internal/adapters/cache"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	clients    map[*websocket.Conn]bool
	currencies chan map[string]float64
	upgrader   websocket.Upgrader
	mu         sync.Mutex
	cache      *cache.RedisCache
}

func NewWebSocketHandler(currencies chan map[string]float64, redis *cache.RedisCache) *WebSocketHandler {
	return &WebSocketHandler{
		clients:    make(map[*websocket.Conn]bool),
		currencies: currencies,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		cache: redis,
	}
}

func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer ws.Close()

	h.mu.Lock()
	h.clients[ws] = true
	h.mu.Unlock()

	log.Println("Client connected")

	// Send a welcome message to the client
	lastBitcoinPrice, err := h.cache.GetBitcoinPrice()

	if err != nil {
		log.Printf("Failed to get last bitcoin price from cache")
		return
	}

	lastUpdatedPrice := map[string]float64{"Bitcoin": lastBitcoinPrice}
	err = ws.WriteJSON(lastUpdatedPrice)
	if err != nil {
		log.Printf("Error sending cacheble bitcoin price to client: %v", err)
		h.mu.Lock()
		delete(h.clients, ws)
		h.mu.Unlock()
		return
	}

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			h.mu.Lock()
			delete(h.clients, ws)
			h.mu.Unlock()
			break
		}
	}
}

func (h *WebSocketHandler) Broadcast() {
	for {
		price, ok := <-h.currencies
		if !ok {
			log.Println("Broadcast channel closed")
			return
		}

		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteJSON(price)
			if err != nil {
				log.Printf("Error sending to client: %v", err)
				client.Close()
				delete(h.clients, client)
			} else {
				log.Println("Price sent to client")
			}
		}
		h.mu.Unlock()
	}
}
