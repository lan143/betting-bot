package roulette

import (
	"errors"
	"fmt"
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

func (c *Controller) PlaceBet(id uuid.UUID, callback *tgbotapi.CallbackQuery, selection string) error {
	if game, ok := c.games[id]; ok {
		err := game.PlaceBet(callback.From.ID, callback.From.UserName, selection)
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
			msg := tgbotapi.NewMessage(game.GetChatId(), fmt.Sprintf("%s successful place bet in %s", callback.From.UserName, selection))
			_, err = c.api.Send(msg)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("game not found")
	}

	return nil
}

func NewController(api *tgbotapi.BotAPI) *Controller {
	return &Controller{api: api}
}
