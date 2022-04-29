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

func (r *Roulette) Generate() Number {
	winIndex := rand.Intn(len(r.numbers))

	return r.numbers[winIndex]
}

func NewRoulette() *Roulette {
	roulette := &Roulette{}
	roulette.Init()

	return roulette
}
