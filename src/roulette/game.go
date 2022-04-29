package roulette

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
)

type Game struct {
	api  *tgbotapi.BotAPI
	step uint8

	id       uuid.UUID
	chatId   int64
	roulette *Roulette
}

func (g *Game) GetId() uuid.UUID {
	return g.id
}

func (g *Game) Init() {
	g.roulette = NewRoulette()
	g.step = 1
}

func (g *Game) Run() {
	for {
		switch g.step {
		case 1:
			g.sendGameMessage()
			break
		}
	}
}

func (g *Game) sendGameMessage() {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup()
	var row = &[]tgbotapi.InlineKeyboardButton{}

	for index, bet := range g.roulette.GetSingleBets(g.id) {
		if index%3 == 0 && len(*row) > 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, *row)
			row = &[]tgbotapi.InlineKeyboardButton{}
		}

		tempRow := append(*row, tgbotapi.NewInlineKeyboardButtonData(bet.Name, bet.Data))
		row = &tempRow
	}

	msg := tgbotapi.NewMessage(g.chatId, "Make your bets")
	msg.ReplyMarkup = keyboard
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}

	g.step = 2
}

func NewGame(api *tgbotapi.BotAPI, id uuid.UUID, chatId int64) *Game {
	return &Game{
		api:    api,
		id:     id,
		chatId: chatId,
	}
}
