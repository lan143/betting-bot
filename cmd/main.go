package main

import (
	"context"
	"main/src"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	container := src.BuildContainer()
	err := container.Invoke(func(a *src.Application) {
		ctx := context.Background()
		err := a.Init(ctx)
		if err != nil {
			panic(err)
		}

		err = a.Run(ctx)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
}
