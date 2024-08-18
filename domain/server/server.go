package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pricing-simulator/domain/constants"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients    map[*websocket.Conn]bool
	currencies chan map[string]float64
	upgrader   websocket.Upgrader
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		currencies: make(chan map[string]float64),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer ws.Close()
	s.mu.Lock()
	s.clients[ws] = true
	s.mu.Unlock()
	log.Println("Client connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			s.mu.Lock()
			delete(s.clients, ws)
			s.mu.Unlock()
			break
		}
		message := string(msg)
		log.Printf("Received: %s", message)
	}
}

func (s *Server) SimulateBitcoinPrice() {
	fmt.Print(constants.BITCOIN)
	price := 50000.0
	for {
		priceChange := rand.Float64()*2000 - 1000
		price += priceChange

		formattedPrice := fmt.Sprintf("%.2f", price)
		log.Printf("Bitcoin new price: %s", formattedPrice)

		var priceFloat float64
		fmt.Sscanf(formattedPrice, "%f", &priceFloat)
		bitcoinMap := make(map[string]float64)
		bitcoinMap[constants.BITCOIN] = priceFloat
		s.currencies <- bitcoinMap

		time.Sleep(60 * time.Second)
	}
}

func (s *Server) Broadcast() {
	for {
		price, ok := <-s.currencies
		if !ok {
			log.Println("Broadcast channel closed")
			return
		}
		log.Printf("Broadcasting price: %.2f", price)

		s.mu.Lock()
		for client := range s.clients {
			err := client.WriteJSON(price)
			if err != nil {
				log.Printf("Error sending to client: %v", err)
				client.Close()
				delete(s.clients, client)
			} else {
				log.Println("Price sent to client")
			}
		}
		s.mu.Unlock()
	}
}
