package logger

import (
	"github.com/pterm/pterm"
	"log/slog"
	"os"
	"strings"
	"time"
)

func New(serviceName string, logLevel string, args ...slog.Attr) *slog.Logger {
	logLevel = strings.ToLower(logLevel)

	baseAttributes := []slog.Attr{
		{
			Key:   "service_name",
			Value: slog.StringValue(serviceName),
		}}

	handler := pterm.NewSlogHandler(&pterm.Logger{
		Formatter:    pterm.LogFormatterColorful,
		Writer:       os.Stdout,
		Level:        getLogLevel(logLevel),
		ShowCaller:   logLevel == "debug",
		CallerOffset: 0,
		ShowTime:     true,
		TimeFormat:   time.RFC3339,
		KeyStyles: map[string]pterm.Style{
			"error":  *pterm.NewStyle(pterm.FgRed, pterm.Bold),
			"err":    *pterm.NewStyle(pterm.FgRed, pterm.Bold),
			"caller": *pterm.NewStyle(pterm.FgGray, pterm.Bold),
		},
		MaxWidth: 80,
	}).WithAttrs(append(baseAttributes, args...))

	return slog.New(handler)
}

func getLogLevel(l string) pterm.LogLevel {
	switch l {
	case "debug":
		return pterm.LogLevelDebug
	case "info":
		return pterm.LogLevelInfo
	case "warn":
		return pterm.LogLevelWarn
	case "error":
		return pterm.LogLevelError
	default:
		return pterm.LogLevelInfo
	}
}
