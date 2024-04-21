package main

import (
	"fmt"

	"github.com/ef4b3f/clog"
)

func main() {
	clog.SetLogLevel(clog.LevelTrace)
	clog.WithTimestamp(true)
	clog.WithCaller(true)
	clog.WithLevelText(true)

	clog.Trace().Msg("hello world!")
	clog.Debug().Msg("hello world!")
	clog.Info().Msg("hello world!")
	clog.Notice().Msg("hello world!")
	clog.Warn().Msg("hello world!")
	clog.Ok().Msg("hello world!")
	clog.Success().Msg("hello world!")
	clog.Error().Msg("hello world!")

	fmt.Println("-------------------------")

	clog.Trace().Any("key", "").Msg("hello world!")
	clog.Debug().Any("key", "str").Msg("hello world!")
	clog.Info().Any("key", true).Msg("hello world!")
	clog.Notice().Any("key", 1).Msg("hello world!")
	clog.Warn().Any("key", nil).Msg("hello world!")
	clog.Ok().Any("key", 0.1).Msg("hello world!")
	clog.Success().Any("key", map[string]string{"foo": "foo", "bar": "bar"}).Msg("hello world!")
	clog.Error().
		Any("key", []string{"foo", "bar"}).
		Any("err", "error message by fields").
		Err(fmt.Errorf("error message by error")).
		Msg("hello world!")

	fmt.Println("-------------------------")

	clog.Fatal().Msg("hello world!")
}
