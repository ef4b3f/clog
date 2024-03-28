package main

import (
	"fmt"

	"github.com/ef4b3f/clog"
)

func main() {
	//clog.SetLogLevel(clog.LevelTrace).WithTimestamp(true).WithCaller(true).WithLevelText(true)
	//clog.SetLogLevel(clog.LevelTrace).WithTimestamp(true).WithCaller(true).WithLevelText(true)
	//clog.SetLogLevel(clog.LevelTrace).WithTimestamp(true).WithCaller(true)
	//clog.SetLogLevel(clog.LevelTrace).WithTimestamp(true)
	clog.SetLogLevel(clog.LevelTrace)
	clog.Trace("hello world!", "trace", 1)
	clog.Debug("hello world!", "debug", "str")
	clog.Notice("hello world!", "notice", true)
	clog.Info("hello world!", "info", []string{"a", "b", "c"})
	clog.Warn("hello world!", "warn", map[string]string{"a": "1", "b": "2"})
	clog.Ok("hello world!")
	clog.Success("hello world!", "success")
	clog.Error("hello world!", "err", fmt.Errorf("err message"))
	clog.Print("hello world!", "a", "1", "b")
	clog.Fatal("hello world!")
}
