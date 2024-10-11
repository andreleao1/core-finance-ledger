package usecase

import (
	"core-finance-ledger/internal/adapters/cache"
	"core-finance-ledger/internal/domain/entity/currencies"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const CURRENCY_NAME = "Bitcoin"

type BitcoinUsecase struct {
	currency currencies.Bitcoin
	cache    *cache.RedisCache
}

func NewBitcoinUsecase(redisCahse *cache.RedisCache) *BitcoinUsecase {
	return &BitcoinUsecase{currency: currencies.NewBitcoin(), cache: redisCahse}
}

func (buc *BitcoinUsecase) SimulateBitcoinPrice() map[string]float64 {
	percentageChange := (rand.Float64()*0.42 - 0.21) / 100
	priceChange := buc.currency.Price * percentageChange
	buc.currency.Price += priceChange

	formattedPrice := fmt.Sprintf("%.2f", buc.currency.Price)
	log.Printf("Bitcoin new price: %s", formattedPrice)

	var priceFloat float64
	fmt.Sscanf(formattedPrice, "%f", &priceFloat)

	return map[string]float64{
		CURRENCY_NAME: priceFloat,
	}
}

func (buc *BitcoinUsecase) StartPriceSimulation(currencies chan map[string]float64) {
	for {
		price := buc.SimulateBitcoinPrice()
		buc.saveToCache(price[buc.currency.CurrencyName])
		currencies <- price
		time.Sleep(60 * time.Second)
	}
}

func (buc *BitcoinUsecase) saveToCache(bitcoinPrice float64) {
	err := buc.cache.SaveBitcoinPrice(bitcoinPrice)

	if err != nil {
		log.Printf("Failed to save bitcoin price in cache")
	}
}
