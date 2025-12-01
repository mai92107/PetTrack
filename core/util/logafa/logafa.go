package logafa

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

type LogLevel slog.Level

var (
	CurrentLevel = slog.LevelDebug
	LogFile      *os.File
)

type LogafaHandler struct {
	level   slog.Leveler
	console *os.File
}

func NewColorHandler(level slog.Leveler, console *os.File) *LogafaHandler {
	return &LogafaHandler{
		level:   level,
		console: console,
	}
}

func (h *LogafaHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

func (h *LogafaHandler) Handle(_ context.Context, r slog.Record) error {
	// timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// get source (file + line)
	file, line := getSource(r.PC)

	// --- 1. 彩色 Console ---
	var colorize *color.Color
	switch r.Level {
	case slog.LevelDebug:
		colorize = color.New(color.FgCyan)
	case slog.LevelInfo:
		colorize = color.New(color.FgGreen)
	case slog.LevelWarn:
		colorize = color.New(color.FgYellow)
	case slog.LevelError:
		colorize = color.New(color.FgRed)
	}

	// message
	msg := r.Message

	// Attributes
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	// Console(彩色)
	consoleLine := fmt.Sprintf("[%s] [%s] [%s:%d] %s%s\n",
		timestamp,
		r.Level.String(),
		file, line,
		colorize.Sprint(msg),
		attrs,
	)
	_, _ = h.console.Write([]byte(consoleLine))

	// --- 3. 寫入檔案（乾淨格式） ---
	if LogFile != nil {
		fileLine := fmt.Sprintf("time=%s level=%s msg=%q %s file=%q line=%d\n",
			time.Now().Format(time.RFC3339),
			r.Level.String(),
			msg,
			attrs,
			file,
			line,
		)
		_, _ = LogFile.Write([]byte(fileLine))
	}

	return nil
}

func (h *LogafaHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *LogafaHandler) WithGroup(name string) slog.Handler {
	return h
}

func Debug(format string, args ...any) {
	slog.Debug(format, args...)
}
func Info(format string, args ...any) {
	slog.Info(format, args...)
}
func Warn(format string, args ...any) {
	slog.Warn(format, args...)
}
func Error(format string, args ...any) {
	slog.Error(format, args...)
}

func getSource(pc uintptr) (file string, line int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", 0
	}
	file, line = fn.FileLine(pc)
	return
}
