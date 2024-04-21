package clog

import (
	"fmt"
	"os"

	"github.com/elliotchance/orderedmap/v2"
)

type Entry struct {
	Logger *Logger
	Level  Level
	Error  error
	Fields *orderedmap.OrderedMap[string, any]
}

func NewEntry(log *Logger) *Entry {
	return &Entry{
		Logger: log,
		Fields: orderedmap.NewOrderedMap[string, any](),
	}
}

func (e *Entry) Any(key string, value any) *Entry {
	e.Fields.Set(key, value)
	return e
}

func (e *Entry) Err(err error) *Entry {
	e.Error = err
	return e
}

func (e *Entry) Msg(format string, args ...any) {
	e.Logger.print(e.Level, fmt.Sprintf(format, args...), e)
	if e.Level == LevelFatal {
		os.Exit(1)
	}
}

func (e *Entry) msg(msg string) {
	e.Logger.print(e.Level, msg, e)
	if e.Level == LevelFatal {
		os.Exit(1)
	}
}
