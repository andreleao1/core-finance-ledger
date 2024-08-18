package main

import (
	"log"
	"net/http"
	"pricing-simulator/domain/server"
)

func main() {
	srv := server.NewServer()
	http.HandleFunc("/currencies", srv.HandleConnections)
	go srv.SimulateBitcoinPrice()
	go srv.Broadcast()

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
