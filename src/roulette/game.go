package roulette

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type UserBet struct {
	UserId    int64
	UserName  string
	Selection string
	Amount    float64
}

type Game struct {
	api  *tgbotapi.BotAPI
	step uint8

	id       uuid.UUID
	chatId   int64
	roulette *Roulette

	bets         []UserBet
	stakeMessage tgbotapi.Message
}

func (g *Game) GetId() uuid.UUID {
	return g.id
}

func (g *Game) GetChatId() int64 {
	return g.chatId
}

func (g *Game) Init() {
	g.roulette = NewRoulette()
	g.step = 1
}

func (g *Game) Run() {
	for {
		switch g.step {
		case 1: // send messages with make bets
			g.sendGameMessage()
			g.step = 2
			break
		case 2: // wait 2 minutes for user bets
			timer := time.NewTimer(20 * time.Second)
			<-timer.C
			msg := tgbotapi.NewEditMessageText(g.chatId, g.stakeMessage.MessageID, "Bets are made, no more bets")
			_, err := g.api.Send(msg)
			if err != nil {
				log.Printf(err.Error())
			}
			g.step = 3
			break
		case 3: // generate win number and settle bets
			g.generateAndSettle()
			g.step = 4
			break
		case 4: // finish
			return
		}
	}
}

func (g *Game) PlaceBet(userId int64, userName string, selection string) error {
	if g.step != 2 {
		return errors.New("time for place bets is out")
	}

	g.bets = append(g.bets, UserBet{
		UserId:    userId,
		UserName:  userName,
		Selection: selection,
		Amount:    5.0,
	})

	return nil
}

func (g *Game) sendGameMessage() {
	var err error
	g.sendSingleSelections()
	g.sendMultipleSelections()

	msg := tgbotapi.NewMessage(g.chatId, "Place your bets gentlemen")
	g.stakeMessage, err = g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}
}

func (g *Game) sendSingleSelections() {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup()
	var row = tgbotapi.NewInlineKeyboardRow()

	for index, selection := range g.roulette.GetSingleSelections(g.id) {
		if index%3 == 0 && len(row) > 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = tgbotapi.NewInlineKeyboardRow()
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(selection.Name, selection.Data))
	}

	if len(row) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	for i, row := range keyboard.InlineKeyboard {
		for j, cell := range row {
			log.Printf("%d-%d: \"%s\" \"%s\"", i, j, cell.Text, *cell.CallbackData)
		}
	}

	msg := tgbotapi.NewMessage(g.chatId, "Make your single bets")
	msg.ReplyMarkup = keyboard
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}
}

func (g *Game) sendMultipleSelections() {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup()
	var row []tgbotapi.InlineKeyboardButton

	for index, selection := range g.roulette.GetMultipleSelections(g.id) {
		if index%3 == 0 && len(row) > 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(selection.Name, selection.Data))
	}

	if len(row) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg := tgbotapi.NewMessage(g.chatId, "Or make your multiple bets")
	msg.ReplyMarkup = keyboard
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}
}

func (g *Game) generateAndSettle() {
	number := g.roulette.Generate()

	msg := tgbotapi.NewMessage(g.chatId, fmt.Sprintf("Roulette is complete! Win: %s %s", number.Num, number.Color))
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}

	message := strings.Builder{}

	for _, bet := range g.bets {
		switch bet.Selection {
		case "Black":
			if number.Color == "Black" {
				message.WriteString(fmt.Sprintf("@%s win %f parrots\n", bet.UserName, bet.Amount*2))
			} else {
				message.WriteString(fmt.Sprintf("@%s lose %f parrots\n", bet.UserName, bet.Amount))
			}
			break
		case "Red":
			if number.Color == "Red" {
				message.WriteString(fmt.Sprintf("@%s win %f parrots\n", bet.UserName, bet.Amount*2))
			} else {
				message.WriteString(fmt.Sprintf("@%s lose %f parrots\n", bet.UserName, bet.Amount))
			}
			break
		default:
			if bet.Selection == number.Num {
				message.WriteString(fmt.Sprintf("@%s win %f parrots\n", bet.UserName, bet.Amount*36))
			} else {
				message.WriteString(fmt.Sprintf("@%s lose %f parrots\n", bet.UserName, bet.Amount))
			}
		}
	}

	text := message.String()
	if len(text) == 0 {
		text = "No bets"
	}

	_, err = g.api.Send(
		tgbotapi.NewMessage(
			g.chatId,
			text,
		),
	)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}
}

func NewGame(api *tgbotapi.BotAPI, id uuid.UUID, chatId int64) *Game {
	return &Game{
		api:    api,
		id:     id,
		chatId: chatId,
	}
}
