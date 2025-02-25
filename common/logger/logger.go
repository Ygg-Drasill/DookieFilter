package logger

import (
	"github.com/pterm/pterm"
	"log/slog"
	"os"
	"time"
)

func New(serviceName string, logLevel string) *slog.Logger {
	handler := pterm.NewSlogHandler(&pterm.Logger{
		Formatter:    pterm.LogFormatterColorful,
		Writer:       os.Stdout,
		Level:        getLogLevel(logLevel),
		ShowCaller:   true,
		CallerOffset: 0,
		ShowTime:     true,
		TimeFormat:   time.RFC3339,
		KeyStyles: map[string]pterm.Style{
			"error":  *pterm.NewStyle(pterm.FgRed, pterm.Bold),
			"err":    *pterm.NewStyle(pterm.FgRed, pterm.Bold),
			"caller": *pterm.NewStyle(pterm.FgGray, pterm.Bold),
		},
		MaxWidth: 80,
	}).WithAttrs([]slog.Attr{
		{
			Key:   "service_name",
			Value: slog.StringValue(serviceName),
		},
	})
	return slog.New(handler)
}

func getLogLevel(l string) pterm.LogLevel {
	switch l {
	case "DISABLE":
		return pterm.LogLevelDisabled
	case "TRACE":
		return pterm.LogLevelTrace
	case "DEBUG":
		return pterm.LogLevelDebug
	case "INFO":
		return pterm.LogLevelInfo
	case "WARN":
		return pterm.LogLevelWarn
	case "ERROR":
		return pterm.LogLevelError
	case "FATAL":
		return pterm.LogLevelFatal
	case "PRINT":
		return pterm.LogLevelPrint
	default:
		return pterm.LogLevelInfo
	}
}
