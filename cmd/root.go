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
	"log"
	"os"

	sub1 "github.com/donovanmods/projectdaedalus-db-tool/cmd/add"
	sub2 "github.com/donovanmods/projectdaedalus-db-tool/cmd/del"
	sub3 "github.com/donovanmods/projectdaedalus-db-tool/cmd/list"
	sub4 "github.com/donovanmods/projectdaedalus-db-tool/cmd/sync"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "pdt",
	Version: "0.1.0",
	Short:   "ProjectDaedalus Database Tool - a CLI utility to manage the Icarus ProjectDaedalus database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		noColor, _ := cmd.Flags().GetBool("no-color")
		verbosity, _ := cmd.Flags().GetCount("verbose")

		viper.Set("color", !noColor)
		viper.Set("verbosity", verbosity)

		logger.SetLogger(verbosity)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	RootCmd.SetVersionTemplate(version())
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pdtconfig.yaml)")
	RootCmd.PersistentFlags().CountP("verbose", "v", "verbose output (may be repeated)")
	RootCmd.PersistentFlags().Bool("dryrun", false, "run without performing any persistent operations")
	RootCmd.PersistentFlags().Bool("color", true, "colorize output")
	RootCmd.PersistentFlags().Bool("no-color", false, "disable color output")

	_ = viper.BindPFlag("dryrun", RootCmd.PersistentFlags().Lookup("dryrun"))

	RootCmd.AddCommand(sub1.AddCmd)
	RootCmd.AddCommand(sub2.DelCmd)
	RootCmd.AddCommand(sub3.ListCmd)
	RootCmd.AddCommand(sub4.SyncCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pdtconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".pdtconfig")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Unable to read config file, this is a required file: ", viper.ConfigFileUsed())
	}
}

func version() string {
	return fmt.Sprintln(RootCmd.Version)
}
