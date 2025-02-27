package logger

import (
	"log/slog"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// Global Log variable
var Log *slog.Logger = slog.Default()

// SetLogger sets the default logger verbosity level and returns a [log/slog.Logger](https://pkg.go.dev/log/slog#Logger)
// verbosity is an integer from 0 to 3, where higher numbers are more verbose
func SetLogger(verbosity int) *slog.Logger {
	if !viper.GetBool("color") {
		pterm.DisableColor()
	}

	// Limit verbosity to 3
	if verbosity > 3 {
		verbosity = 3
	}

	// Create a new slog handler with the default PTerm logger
	handler := pterm.NewSlogHandler(&pterm.DefaultLogger)

	// Create a new slog logger with the handler
	logger := slog.New(handler)

	var level pterm.LogLevel
	switch verbosity {
	case 3:
		level = pterm.LogLevelDebug
	case 2:
		level = pterm.LogLevelInfo
	case 1:
		level = pterm.LogLevelWarn
	default:
		level = pterm.LogLevelError
	}

	// Change the log level to debug to enable debug messages
	pterm.DefaultLogger.Level = level

	Log = logger

	return logger
}

func TestLogger() {
	Log.Debug("This is a debug message")
	Log.Info("This is an info message")
	Log.Warn("This is a warning message")
	Log.Error("This is an error message")
}
