package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan float64
	upgrader  websocket.Upgrader
	mu        sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan float64),
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
	price := 50000.0
	for {
		priceChange := rand.Float64()*2000 - 1000
		price += priceChange

		formattedPrice := fmt.Sprintf("%.2f", price)
		log.Printf("Broadcasting new price: %s", formattedPrice)

		var priceFloat float64
		fmt.Sscanf(formattedPrice, "%f", &priceFloat)
		s.broadcast <- priceFloat

		time.Sleep(60 * time.Second)
	}
}

func (s *Server) Broadcast() {
	for {
		price, ok := <-s.broadcast
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
