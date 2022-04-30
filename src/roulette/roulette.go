package roulette

import (
	"github.com/google/uuid"
	"math/rand"
)

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

type Selection struct {
	Name string
	Data string
}

type Roulette struct {
	numbers []Number

	winNumber *Number
}

func (r *Roulette) GetSingleSelections(id uuid.UUID) []Selection {
	var selections []Selection

	for _, number := range r.numbers {
		selections = append(selections, Selection{
			Name: number.Num + " - " + string(number.Color),
			Data: "roulette-selections;" + id.String() + ";" + number.Num,
		})
	}

	return selections
}

func (r *Roulette) GetMultipleSelections(id uuid.UUID) []Selection {
	var selections = []Selection{
		{
			Name: "Black",
			Data: "roulette-selections;" + id.String() + ";Black",
		},
		{
			Name: "Red",
			Data: "roulette-selections;" + id.String() + ";Red",
		},
	}

	return selections
}

func (r *Roulette) EndRoulette() Number {
	if r.winNumber != nil {
		return *r.winNumber
	}

	winIndex := rand.Intn(len(r.numbers))
	r.winNumber = &r.numbers[winIndex]

	return *r.winNumber
}

func (r *Roulette) GetBetWin(bet *UserBet) uint {
	switch bet.Selection {
	case "Black":
		if r.winNumber.Color == Black {
			return bet.Amount * 2
		} else {
			return 0
		}
	case "Red":
		if r.winNumber.Color == Red {
			return bet.Amount * 2
		} else {
			return 0
		}
	default:
		if bet.Selection == r.winNumber.Num {
			return bet.Amount * 36
		} else {
			return 0
		}
	}
}

func NewRoulette() *Roulette {
	roulette := &Roulette{}
	roulette.Init()

	return roulette
}
