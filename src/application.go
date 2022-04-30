package src

import (
	"context"
	"log"
	"main/src/config"
	"main/src/telegram"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	config *config.Config
	bot    *telegram.Bot

	wg   sync.WaitGroup
	sigs chan os.Signal
}

func (a *Application) Init(ctx context.Context) error {
	log.Printf("application: init")

	a.sigs = make(chan os.Signal, 1)
	signal.Notify(a.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	err := a.config.Init(ctx)
	if err != nil {
		return err
	}

	err = a.bot.Init(ctx, &a.wg, a.config.GetBotToken())
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Run(ctx context.Context) error {
	log.Printf("application.run: start")

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go a.processSignals(cancelFunc)

	err := a.bot.Run(cancelCtx)
	if err != nil {
		return err
	}

	log.Println("application.run: running")

	a.wg.Wait()

	log.Println("application: graceful shutdown.")

	return nil
}

func (a *Application) processSignals(cancelFun context.CancelFunc) {
	select {
	case <-a.sigs:
		log.Println("application: received shutdown signal from OS")
		cancelFun()
		break
	}
}

func NewApplication(
	config *config.Config,
	bot *telegram.Bot,
) *Application {
	return &Application{
		config: config,
		bot:    bot,
	}
}
