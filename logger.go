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

func (l Level) String() string {
	return levelText[l]
}

const (
	_ Level = iota - 1
	LevelTrace
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelOk
	LevelSuccess
	LevelError
	LevelFatal
)

type LevelStyle struct {
	Color   lipgloss.Color
	Icon    lipgloss.Style
	Text    lipgloss.Style
	Message lipgloss.Style
	Key     lipgloss.Style
}

var (
	levelText = [...]string{
		LevelTrace:   "trace",
		LevelDebug:   "debug",
		LevelInfo:    "info",
		LevelNotice:  "notice",
		LevelWarn:    "warn",
		LevelOk:      "ok",
		LevelSuccess: "success",
		LevelError:   "error",
		LevelFatal:   "fatal",
	}

	Styles = [...]LevelStyle{
		LevelTrace: {
			Color:   lipgloss.Color("51"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelDebug: {
			Color:   lipgloss.Color("102"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelInfo: {
			Color:   lipgloss.Color(""),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelNotice: {
			Color:   lipgloss.Color("15"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("•"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle().Bold(true),
		},
		LevelWarn: {
			Color:   lipgloss.Color("214"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("⚠"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelOk: {
			Color:   lipgloss.Color("27"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✔"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelSuccess: {
			Color:   lipgloss.Color("47"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✔"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelError: {
			Color:   lipgloss.Color("196"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✖"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle(),
		},
		LevelFatal: {
			Color:   lipgloss.Color("160"),
			Icon:    lipgloss.NewStyle().Bold(true).SetString("✖"),
			Text:    lipgloss.NewStyle().Bold(true),
			Key:     lipgloss.NewStyle().Bold(true).Faint(true),
			Message: lipgloss.NewStyle().Bold(true),
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
		Level:         LevelInfo,
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

func (l *Logger) renderTimestamp() string {
	if !l.ShowTime {
		return ""
	}
	return fmt.Sprintf("%s %s ",
		lipgloss.NewStyle().Foreground(gray).Render(time.Now().Format(l.TimeFormat)),
		divide.Render(),
	)
}

func (l *Logger) renderLevelText(level Level) string {
	if !l.ShowLevelText {
		return ""
	}
	text, style := strings.ToUpper(level.String()), Styles[level]
	return fmt.Sprintf("%s%s ",
		style.Text.Foreground(style.Color).Width(8).SetString(text).Render(),
		divide.Render(),
	)
}

func (l *Logger) print(level Level, msg string, e *Entry) {
	if level < l.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	style := Styles[level]
	timestamp := l.renderTimestamp()

	_, _ = fmt.Fprint(l.Writer,
		timestamp,
		style.Icon.Foreground(style.Color).Render(""),
		l.renderLevelText(level),
		style.Message.Render(msg),
	)

	if e.Error != nil {
		e.Any("err", e.Error)
	}
	if l.ShowCaller {
		path, line := l.getCallerInfo()
		e.Any("caller", lipgloss.NewStyle().Foreground(gray).Render(fmt.Sprintf("%s:%d", path, line)))
	}
	i := 0
	for it := e.Fields.Front(); it != nil; it = it.Next() {
		key, value := it.Key, it.Value
		i++
		keyStyle := style.Key.Copy().Foreground(style.Color)
		if key == "caller" {
			keyStyle.Foreground(gray)
		}
		if value == nil {
			value = ""
		}
		value = fmt.Sprint(value)
		if key != "" && value != "" {
			key += ": "
		}
		argPrefix := "├─"
		if i == e.Fields.Len() {
			argPrefix = "└─"
		}
		_, _ = fmt.Fprintf(l.Writer,
			"\n  %s %s%s",
			lipgloss.NewStyle().Foreground(gray).Faint(true).Render(argPrefix),
			keyStyle.Render(key),
			value,
		)
	}

	_, _ = fmt.Fprintln(l.Writer)
}

func (l *Logger) newEntry(level Level) *Entry {
	e := NewEntry(l)
	e.Level = level
	return e
}

func (l *Logger) Trace() *Entry {
	return l.newEntry(LevelTrace)
}

func (l *Logger) Debug() *Entry {
	return l.newEntry(LevelDebug)
}

func (l *Logger) Info() *Entry {
	return l.newEntry(LevelInfo)
}

func (l *Logger) Notice() *Entry {
	return l.newEntry(LevelNotice)
}

func (l *Logger) Warn() *Entry {
	return l.newEntry(LevelWarn)
}

func (l *Logger) Ok() *Entry {
	return l.newEntry(LevelOk)
}

func (l *Logger) Success() *Entry {
	return l.newEntry(LevelSuccess)
}

func (l *Logger) Error() *Entry {
	return l.newEntry(LevelError)
}

func (l *Logger) Fatal() *Entry {
	return l.newEntry(LevelFatal)
}

func Trace() *Entry {
	return logger.Trace()
}

func Debug() *Entry {
	return logger.Debug()
}

func Info() *Entry {
	return logger.Info()
}

func Notice() *Entry {
	return logger.Notice()
}

func Warn() *Entry {
	return logger.Warn()
}

func Ok() *Entry {
	return logger.Ok()
}

func Success() *Entry {
	return logger.Success()
}

func Error() *Entry {
	return logger.Error()
}

func Fatal() *Entry {
	return logger.Fatal()
}
