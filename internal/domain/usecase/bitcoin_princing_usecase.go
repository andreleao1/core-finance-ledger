package usecase

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const CURRENCY_NAME = "Bitcoin"

type PriceSimulator struct {
	price float64
}

func NewPriceSimulator() *PriceSimulator {
	return &PriceSimulator{
		price: 339687.03, // Preço inicial do Bitcoin
	}
}

// Simula a variação do preço do Bitcoin entre -0,21% e +0,21%
func (ps *PriceSimulator) SimulateBitcoinPrice() map[string]float64 {
	// Variação percentual entre -0,21% e +0,21%
	percentageChange := (rand.Float64()*0.42 - 0.21) / 100
	priceChange := ps.price * percentageChange
	ps.price += priceChange

	formattedPrice := fmt.Sprintf("%.2f", ps.price)
	log.Printf("Bitcoin new price: %s", formattedPrice)

	var priceFloat float64
	fmt.Sscanf(formattedPrice, "%f", &priceFloat)

	return map[string]float64{
		CURRENCY_NAME: priceFloat,
	}
}

func (ps *PriceSimulator) StartPriceSimulation(currencies chan map[string]float64) {
	for {
		price := ps.SimulateBitcoinPrice()
		currencies <- price
		time.Sleep(60 * time.Second) // Simula a cada 60 segundos
	}
}
