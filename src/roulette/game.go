package roulette

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"log"
	"main/src/entities"
	"main/src/repositories"
	"strings"
	"time"
)

type UserBet struct {
	User      *entities.User
	Selection string
	Amount    uint
}

type Game struct {
	api             *tgbotapi.BotAPI
	usersRepository repositories.UsersRepository

	id       uuid.UUID
	chatId   int64
	step     uint8
	roulette *Roulette

	bets                       []UserBet
	messagesWithSelections     []int
	messageWithInformation     int
	messageWithInformationText string
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
			g.messagesWithSelections = g.sendGameMessage()

			g.messageWithInformationText = "Place your bets gentlemen\n\n"
			msg := tgbotapi.NewMessage(g.chatId, g.messageWithInformationText)
			send, err := g.api.Send(msg)
			if err != nil {
				log.Printf("roulette.send-game-message.send: %s", err.Error())
				return
			} else {
				g.messageWithInformation = send.MessageID
			}

			g.step = 2
			break
		case 2: // wait 20 seconds for user bets
			timer := time.NewTimer(20 * time.Second)
			<-timer.C
			for _, id := range g.messagesWithSelections {
				msg := tgbotapi.NewDeleteMessage(g.chatId, id)
				_, err := g.api.Send(msg)
				if err != nil {
					log.Printf(err.Error())
				}
			}

			g.step = 3
			break
		case 3: // generate win number and settle bets
			msg := tgbotapi.NewEditMessageText(
				g.GetChatId(),
				g.messageWithInformation,
				"Bets are made, no more bets. Roulette ends up moving...",
			)
			_, err := g.api.Send(msg)
			if err != nil {
				log.Printf("game.place-bet.send: %s", err.Error())
			}

			timer := time.NewTimer(3 * time.Second)
			<-timer.C
			g.generateAndSettle()
			g.step = 4
			break
		case 4: // finish
			return
		}
	}
}

func (g *Game) PlaceBet(userId int64, userName string, selection string, stake float64) error {
	if g.step != 2 {
		return errors.New("time for place bets is out")
	}

	user, err := g.usersRepository.GetUserByExternalId(userId, g.chatId)
	if err != nil {
		log.Printf("game.place-bet.get-user: %s", err.Error())
		return errors.New("bot error: can't place bet")
	}

	if user == nil {
		user = &entities.User{
			Id:             uuid.New(),
			ExternalId:     userId,
			ExternalChatId: g.chatId,
			UserName:       userName,
			Balance:        100 * 100,
		}
		err := g.usersRepository.Save(user)
		if err != nil {
			log.Printf("game.place-bet.save-user: %s", err.Error())
			return errors.New("bot error: can't place bet")
		}
	}

	stakeInCents := uint(stake * 100)

	if user.Balance < stakeInCents {
		return errors.New("not enough money")
	}

	err = g.usersRepository.UpdateBalance(user, -stakeInCents)
	if err != nil {
		log.Printf("game.place-bet.update-balance: %s", err.Error())
		return errors.New("bot error: can't place bet")
	}

	g.bets = append(g.bets, UserBet{
		User:      user,
		Selection: selection,
		Amount:    stakeInCents,
	})

	g.messageWithInformationText = g.messageWithInformationText + fmt.Sprintf(
		"%s placed bet on %s. Balance: %.2f parrots\n",
		user.UserName,
		selection,
		float64(user.Balance)/100,
	)

	msg := tgbotapi.NewEditMessageText(
		g.GetChatId(),
		g.messageWithInformation,
		g.messageWithInformationText,
	)
	_, err = g.api.Send(msg)
	if err != nil {
		log.Printf("game.place-bet.send: %s", err.Error())
	}

	return nil
}

func (g *Game) sendGameMessage() []int {
	var messageIds []int
	messageIds = append(messageIds, g.sendSingleSelections(), g.sendMultipleSelections())

	return messageIds
}

func (g *Game) sendSingleSelections() int {
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
	send, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}

	return send.MessageID
}

func (g *Game) sendMultipleSelections() int {
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
	send, err := g.api.Send(msg)
	if err != nil {
		log.Printf("roulette.send-game-message.send: %s", err.Error())
	}

	return send.MessageID
}

func (g *Game) generateAndSettle() {
	number := g.roulette.EndRoulette()

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("Roulette is complete! Victory: %s %s\n\n", number.Num, number.Color))

	for _, bet := range g.bets {
		winAmount := g.roulette.GetBetWin(&bet)
		if winAmount > 0 {
			err := g.usersRepository.UpdateBalance(bet.User, bet.Amount+winAmount)
			if err != nil {
				log.Printf("game.generate-and-settle.update-balance: %s", err.Error())
				break
			}

			message.WriteString(
				fmt.Sprintf(
					"✌️ @%s bet on %s and win %.2f parrots. Balance: %.2f parrots.\n",
					bet.User.UserName,
					bet.Selection,
					float64(winAmount)/100,
					float64(bet.User.Balance)/100,
				),
			)
		} else {
			message.WriteString(
				fmt.Sprintf(
					"😥️ @%s bet on %s and lose %.2f parrots. Balance: %.2f parrots.\n",
					bet.User.UserName,
					bet.Selection,
					float64(bet.Amount)/100,
					float64(bet.User.Balance)/100,
				),
			)
		}
	}

	msg := tgbotapi.NewEditMessageText(
		g.GetChatId(),
		g.messageWithInformation,
		message.String(),
	)
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("game.place-bet.send: %s", err.Error())
	}
}

func NewGame(
	api *tgbotapi.BotAPI,
	usersRepository repositories.UsersRepository,
	id uuid.UUID,
	chatId int64,
) *Game {
	return &Game{
		api:             api,
		usersRepository: usersRepository,
		id:              id,
		chatId:          chatId,
	}
}
