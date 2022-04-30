package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
	"main/src/repositories"
	"main/src/roulette"
	"strings"
	"sync"
)

type Bot struct {
	api *tgbotapi.BotAPI
	wg  *sync.WaitGroup

	rouletteController *roulette.Controller
	usersRepository    repositories.UsersRepository

	isShutdown bool
}

func (b *Bot) Init(ctx context.Context, wg *sync.WaitGroup, token string) error {
	log.Printf("bot: init")

	var err error
	b.wg = wg
	b.api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	b.rouletteController = roulette.NewController(b.api, b.usersRepository)
	b.rouletteController.Init()

	log.Printf("bot.init: authorized on account %s", b.api.Self.UserName)
	log.Printf("bot.init: complete")

	return nil
}

func (b *Bot) Run(ctx context.Context) error {
	log.Printf("bot.run: running")
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := b.api.GetUpdatesChan(u)

		for {
			if b.isShutdown {
				break
			}

			select {
			case <-ctx.Done():
				log.Printf("bot: shutdown")
				b.isShutdown = true
				break
			case update := <-updates:
				go b.handleUpdate(update)
				break
			}
		}
	}()

	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		b.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update.CallbackQuery)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	var text string
	if message.IsCommand() {
		text = message.Command()
	} else {
		text = message.Text
	}

	log.Printf("[%s] %+v", message.From.UserName, text)

	if message.IsCommand() {
		switch message.Command() {
		case "new_roulette":
			b.newRoulette(message.Chat.ID)
			break
		default:
			msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know that command")
			msg.ReplyToMessageID = message.MessageID

			_, err := b.api.Send(msg)
			if err != nil {
				log.Printf("bot.handle-update.send: %s", err.Error())
			}
		}
	}
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	log.Printf("[%s] %+v", callback.From.UserName, callback.Data)

	data := strings.Split(callback.Data, ";")
	log.Printf("data: %+v", data)

	if len(data) > 0 {
		switch data[0] {
		case "roulette-selections":
			id, err := uuid.Parse(data[1])
			if err != nil {
				log.Printf(err.Error())
			} else {
				b.rouletteController.PlaceBet(id, callback, data[2])
			}
			break
		default:
			log.Printf("unsupported %s", data[0])
		}
	}
}

func (b *Bot) newRoulette(chatId int64) {
	b.rouletteController.NewGame(chatId)
}

func NewBot(usersRepository repositories.UsersRepository) *Bot {
	return &Bot{usersRepository: usersRepository}
}
