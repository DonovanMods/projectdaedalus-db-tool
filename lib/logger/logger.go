package logger

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// Global Log variable
var Log *pterm.Logger

// SetLogger sets the default logger verbosity level and returns a logger instance
func SetLogger(verbosity int) *pterm.Logger {
	if !viper.GetBool("color") {
		pterm.DisableColor()
	}

	Log = pterm.DefaultLogger.WithLevel(getLevel(verbosity))

	return Log
}

// Helper Functions
func Panic(err error) {
	panic(err)
}

func Fatal(err error) {
	Log.Fatal(err.Error())
	os.Exit(1)
}

func Error(msg string) {
	Log.Error(msg)
}

func Warn(msg string) {
	Log.Warn(msg)
}

func Info(msg string) {
	Log.Info(msg)
}

func Debug(msg string) {
	Log.Debug(msg)
}

func Trace(msg string) {
	Log.Trace(msg)
}

// Private functions

func getLevel(verbosity int) pterm.LogLevel {
	if verbosity >= 4 {
		return pterm.LogLevelTrace
	}

	switch verbosity {
	case 3:
		return pterm.LogLevelDebug
	case 2:
		return pterm.LogLevelInfo
	case 1:
		return pterm.LogLevelWarn
	default:
		return pterm.LogLevelError
	}
}

func init() {
	if Log == nil {
		SetLogger(0)
	}
}
