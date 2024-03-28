package clog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	gray   = lipgloss.Color("240")
	divide = lipgloss.NewStyle().Foreground(gray).Faint(true).SetString("∣")
)

type Level int

const (
	LevelNone Level = iota - 1
	LevelTrace
	LevelDebug
	LevelNotice
	LevelInfo
	LevelWarn
	LevelOk
	LevelSuccess
	LevelError
	LevelFatal
	LevelPrint
)

type LevelStyle struct {
	Color   lipgloss.Color
	Icon    lipgloss.Style
	Text    lipgloss.Style
	Message lipgloss.Style
	Key     lipgloss.Style
}

var (
	Styles = [...]LevelStyle{
		LevelTrace: {
			Color:   lipgloss.Color("63"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("TRACE"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelDebug: {
			Color:   lipgloss.Color("145"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("DEBUG"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelNotice: {
			Color:   lipgloss.Color("192"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("NOTICE"),
			Key:     lipgloss.NewStyle().Bold(true),
			Message: lipgloss.NewStyle(),
		},
		LevelInfo: {
			Color:   lipgloss.Color("86"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("INFO"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelWarn: {
			Color:   lipgloss.Color("3"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("⚠"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("WARN"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelOk: {
			Color:   lipgloss.Color("33"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✔"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("OK"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelSuccess: {
			Color:   lipgloss.Color("34"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✔"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("SUCCESS"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelError: {
			Color:   lipgloss.Color("204"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✖"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("ERROR"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelFatal: {
			Color:   lipgloss.Color("134"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✖"),
			Text:    lipgloss.NewStyle().Bold(true).SetString("FATAL"),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelPrint: {
			Color:   lipgloss.Color(""),
			Icon:    lipgloss.NewStyle(),
			Text:    lipgloss.NewStyle(),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
	}
)

type Argument struct {
	Key   string
	Value any
}

type Logger struct {
	mu            sync.Mutex
	Writer        io.Writer
	Level         Level
	ShowLevelText bool
	ShowCaller    bool
	ShowTime      bool
	TimeFormat    string
}

var logger = New()

func New() *Logger {
	return &Logger{
		Writer:        os.Stderr,
		Level:         LevelNotice,
		ShowLevelText: false,
		ShowCaller:    false,
		ShowTime:      false,
		TimeFormat:    "2006-01-02 15:04:05",
	}
}

func WithLevelText(with bool) *Logger {
	return logger.WithLevelText(with)
}

func WithCaller(with bool) *Logger {
	return logger.WithCaller(with)
}

func WithTimestamp(with bool) *Logger {
	return logger.WithTimestamp(with)
}

func SetTimeFormat(timeFormat string) *Logger {
	return logger.SetTimeFormat(timeFormat)
}

func SetWriter(writer io.Writer) *Logger {
	return logger.SetWriter(writer)
}

func SetLogLevel(level Level) *Logger {
	return logger.SetLogLevel(level)
}

func (l *Logger) WithLevelText(with bool) *Logger {
	l.ShowLevelText = with
	return l
}

func (l *Logger) WithCaller(with bool) *Logger {
	l.ShowCaller = with
	return l
}

func (l *Logger) WithTimestamp(with bool) *Logger {
	l.ShowTime = with
	return l
}

func (l *Logger) SetTimeFormat(timeFormat string) *Logger {
	l.TimeFormat = timeFormat
	return l
}

func (l *Logger) SetWriter(writer io.Writer) *Logger {
	l.Writer = writer
	return l
}

func (l *Logger) SetLogLevel(level Level) *Logger {
	l.Level = level
	return l
}

func (l *Logger) combineArgs(args ...any) []Argument {
	return Args(args...)
}

func (l *Logger) getCallerInfo() (path string, line int) {
	if !l.ShowCaller {
		return
	}

	_, path, line, _ = runtime.Caller(4)
	_, callerBase, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(callerBase)
	basepath = strings.ReplaceAll(basepath, "\\", "/")

	path = strings.TrimPrefix(path, basepath)

	return
}

func (l *Logger) renderTimestamp(level Level) string {
	if !l.ShowTime || level == LevelPrint {
		return ""
	}
	return fmt.Sprintf("%s %s ",
		lipgloss.NewStyle().Foreground(gray).Render(time.Now().Format(l.TimeFormat)),
		divide.Render(),
	)
}

func (l *Logger) renderLevelText(level Level) string {
	if !l.ShowLevelText || level == LevelPrint {
		return ""
	}
	style := Styles[level]
	return fmt.Sprintf("%s%s ",
		style.Text.Foreground(style.Color).Width(8).Render(),
		divide.Render(),
	)
}

func (l *Logger) print(level Level, msg string, args []Argument) {
	if level < l.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	style := Styles[level]
	timestamp := l.renderTimestamp(level)

	_, _ = fmt.Fprint(l.Writer,
		style.Icon.Foreground(style.Color).Render(""),
		l.renderLevelText(level),
		timestamp,
		style.Message.Render(msg),
	)

	if l.ShowCaller {
		path, line := l.getCallerInfo()
		args = append(args, Argument{
			Key:   "caller",
			Value: lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%s:%d", path, line)),
		})
	}
	for i, argument := range args {
		keyStyle := style.Key.Copy().Foreground(style.Color)
		key, value := argument.Key, ""
		if key == "caller" {
			keyStyle.Foreground(gray)
		}
		if argument.Value != nil {
			value = fmt.Sprint(argument.Value)
		}
		if key != "" && value != "" {
			key += ": "
		}
		argPrefix := "├─"
		if i+1 == len(args) {
			argPrefix = "└─"
		}
		_, _ = fmt.Fprintf(l.Writer,
			"\n  %s %s%s",
			lipgloss.NewStyle().Foreground(gray).Render(argPrefix),
			keyStyle.Render(key),
			value,
		)
	}

	_, _ = fmt.Fprintln(l.Writer)
}

func Args(args ...any) []Argument {
	var result []Argument
	for i := 0; i < len(args); i += 2 {
		argument := Argument{Key: fmt.Sprint(args[i])}
		if i+1 < len(args) {
			argument.Value = args[i+1]
		}
		result = append(result, argument)
	}
	return result
}

func (l *Logger) Trace(msg string, args ...any) {
	l.print(LevelTrace, msg, l.combineArgs(args...))
}

func (l *Logger) Debug(msg string, args ...any) {
	l.print(LevelDebug, msg, l.combineArgs(args...))
}

func (l *Logger) Notice(msg string, args ...any) {
	l.print(LevelNotice, msg, l.combineArgs(args...))
}

func (l *Logger) Info(msg string, args ...any) {
	l.print(LevelInfo, msg, l.combineArgs(args...))
}

func (l *Logger) Warn(msg string, args ...any) {
	l.print(LevelWarn, msg, l.combineArgs(args...))
}

func (l *Logger) Ok(msg string, args ...any) {
	l.print(LevelOk, msg, l.combineArgs(args...))
}

func (l *Logger) Success(msg string, args ...any) {
	l.print(LevelSuccess, msg, l.combineArgs(args...))
}

func (l *Logger) Error(msg string, args ...any) {
	l.print(LevelError, msg, l.combineArgs(args...))
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.print(LevelFatal, msg, l.combineArgs(args...))
	os.Exit(1)
}

func (l *Logger) Print(msg string, args ...any) {
	l.print(LevelPrint, msg, l.combineArgs(args...))
}

func Trace(msg string, args ...any) {
	logger.Trace(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Notice(msg string, args ...any) {
	logger.Notice(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Ok(msg string, args ...any) {
	logger.Ok(msg, args...)
}

func Success(msg string, args ...any) {
	logger.Success(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	logger.Fatal(msg, args...)
}

func Print(msg string, args ...any) {
	logger.Print(msg, args...)
}
