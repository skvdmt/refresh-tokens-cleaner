package main

import (
	"os"

	"github.com/skvdmt/refresh-tokens-cleaner/internal"
	"github.com/skvdmt/refresh-tokens-cleaner/internal/model"
)

func main() {
	var err error
	if model.Log, err = model.NewLogger(); err != nil {
		panic(err)
	}
	defer func() {
		if err := model.Log.Close(); err != nil {
			panic(err)
		}
	}()
	if model.Conf, err = model.NewConfig(); err != nil {
		model.Log.Error.Error(err.Error())
		os.Exit(1)
	}
	if err := internal.RemoveExpiredRefreshTokens(); err != nil {
		model.Log.Error.Error(err.Error())
		os.Exit(1)
	}
}
