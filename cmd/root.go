/*
Copyright © 2025 Donovan C. Young <dyoung522@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "imt",
	Version: "0.1.0",
	Short:   "CLI tool to help manage the Icarus ProjectDaedalus database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		noColor, _ := cmd.Flags().GetBool("no-color")
		if noColor {
			color.NoColor = true
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetVersionTemplate(version())
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.imt.yaml)")
	rootCmd.PersistentFlags().CountP("verbose", "v", "verbose output (may be repeated)")
	rootCmd.PersistentFlags().Bool("dryrun", false, "run without performing any persistent operations")
	rootCmd.PersistentFlags().Bool("color", true, "colorize output")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")

	_ = viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	verbosity, _ := rootCmd.Flags().GetCount("verbose")
	viper.Set("verbosity", verbosity)

	setLogger(verbosity)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".imtconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".imtconfig")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug(fmt.Sprint("Using config file", viper.ConfigFileUsed()))
	}
}

func version() string {
	return fmt.Sprintln(rootCmd.Version)
}

func setLogger(verbosity int) {
	var level slog.Level

	if verbosity > 3 {
		verbosity = 3
	}

	switch verbosity {
	case 3:
		level = slog.LevelDebug
	case 2:
		level = slog.LevelInfo
	case 1:
		level = slog.LevelWarn
	default:
		level = slog.LevelError
	}

	slog.SetLogLoggerLevel(level)

	// logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	// slog.SetDefault(logger)
}
