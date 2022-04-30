package roulette

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"main/src/repositories"
)

type Controller struct {
	api             *tgbotapi.BotAPI
	usersRepository repositories.UsersRepository

	games map[uuid.UUID]*Game
}

func (c *Controller) Init() {
	c.games = make(map[uuid.UUID]*Game, 100)
}

func (c *Controller) NewGame(chatId int64) {
	game := NewGame(c.api, c.usersRepository, uuid.New(), chatId)
	c.games[game.GetId()] = game

	game.Init()
	go game.Run()
}

func (c *Controller) PlaceBet(id uuid.UUID, callback *tgbotapi.CallbackQuery, selection string) error {
	if game, ok := c.games[id]; ok {
		err := game.PlaceBet(callback.From.ID, callback.From.UserName, selection, 5.0)
		if err != nil {
			msg := tgbotapi.CallbackConfig{
				CallbackQueryID: callback.ID,
				Text:            err.Error(),
				ShowAlert:       true,
			}
			_, err = c.api.Send(msg)
			if err != nil {
				return err
			}
		} else {

		}
	} else {
		return errors.New("game not found")
	}

	return nil
}

func NewController(api *tgbotapi.BotAPI, usersRepository repositories.UsersRepository) *Controller {
	return &Controller{api: api, usersRepository: usersRepository}
}
