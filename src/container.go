package src

import (
	"go.uber.org/dig"
	"main/src/config"
	"main/src/repositories/memory"
	"main/src/telegram"
)

func BuildContainer() *dig.Container {
	container := dig.New()
	processError(container.Provide(NewApplication))

	// Config
	processError(container.Provide(config.NewConfig))

	// Repositories
	processError(container.Provide(memory.NewInMemoryUsersRepository))

	// Telegram
	processError(container.Provide(telegram.NewBot))

	return container
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}
