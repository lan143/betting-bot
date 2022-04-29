package roulette

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type Controller struct {
	api *tgbotapi.BotAPI

	games map[uuid.UUID]*Game
}

func (c *Controller) Init() {
	c.games = make(map[uuid.UUID]*Game, 100)
}

func (c *Controller) NewGame(chatId int64) {
	game := NewGame(c.api, uuid.New(), chatId)
	c.games[game.GetId()] = game

	game.Init()
	go game.Run()
}

func NewController(api *tgbotapi.BotAPI) *Controller {
	return &Controller{api: api}
}
