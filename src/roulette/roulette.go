package roulette

import "github.com/google/uuid"

type Color string

const (
	Green Color = "Green"
	Black       = "Black"
	Red         = "Red"
)

type Number struct {
	Num   string
	Color Color
}

type Bet struct {
	Name string
	Data string
}

type Roulette struct {
	numbers []Number
}

func (r *Roulette) GetSingleBets(id uuid.UUID) []Bet {
	var bets []Bet

	for _, number := range r.numbers {
		bets = append(bets, Bet{
			Name: number.Num + " - " + string(number.Color),
			Data: "roulette-bets;" + id.String() + ";" + number.Num + "-" + string(number.Color),
		})
	}

	return bets
}

func NewRoulette() *Roulette {
	roulette := &Roulette{}
	roulette.Init()

	return roulette
}
