package currencies

type Bitcoin struct {
	CurrencyName string
	Price        float64
}

func NewBitcoin() Bitcoin {
	return Bitcoin{
		CurrencyName: "Bitcoin",
		Price:        339687.03,
	}
}
