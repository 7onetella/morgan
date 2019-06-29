package cmd

// MIT License

// Copyright (c) 2019 7onetella

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logVal string
var logLevel int

const (
	// DEBUG is debug log level
	DEBUG = 2
	// INFO is info log level
	INFO = 1
	// NONE is no logging
	NONE = 0
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "morgan",
	Short: "morgan is a suite of chatops CLI commands",
	Long:  `morgan is a suite of chatops CLI commands`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		switch logVal {
		case "debug":
			logLevel = DEBUG
		case "info":
			logLevel = INFO
		default:
			logLevel = 0
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	UpdateExampleOnChildren(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.morgan.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	term, ok := os.LookupEnv("TERM_PROGRAM")
	if ok && term == "iTerm.app" {
		_isTerminal = true
	}

	pflags.StringVar(&logVal, "log", "", "logging level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mvk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".morgan")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var indentation = "  "
var checkmark = "\u2713 "
var xmark = "\u2717 "
var bullet = "\u2022 "
var _isTerminal bool
var loggingLevel = 0

// CheckArgs checks to see if minimum number of args are provided
func CheckArgs(args []string, minimum int) {
	if len(args) < minimum {
		Println(red(indentation + xmark + "missing arguments"))
		os.Exit(1)
	}
}

// ExitOnError exits when error occurs
func ExitOnError(err error, action string) {
	if err != nil {
		Println(red(indentation + xmark + action))
		if logLevel == DEBUG {
			Print("\n")
			IndentRed(err.Error())
		}
		os.Exit(1)
	}

	// if there is no error but debug is set
	// this may not sound right but we need to show what actions were successful
	// before showing the error from the action that failed at the end
	if logLevel >= DEBUG {
		Println(green(indentation + checkmark + action))
	}
}

// ExitOn exits when error occurs
func ExitOn(err error) {
	if err != nil {
		// if logLevel == DEBUG {
		Print(red(indentation + xmark + err.Error()))
		Print("\n")
		// IndentRed(err.Error())
		// }
		os.Exit(1)
	}
}

// ExitOnErrorWithDetail exits when error occurs, prints addtional info
func ExitOnErrorWithDetail(err error, action, stdout, errout string) {
	if err != nil {
		Println(red(indentation + xmark + action))
		if logLevel == DEBUG {
			Print("\n")
			IndentRed(stdout)
			IndentRed(errout)
		}
		os.Exit(1)
	}

	// if there is no error but debug is set
	// this may not sound right but we need to show what actions were successful
	// before showing the error from the action that failed at the end
	if logLevel >= DEBUG {
		Println(green(indentation + checkmark + action))
	}
}

// Success prints green with checkmark preceding the message
func Success(msg string) {
	if logLevel != DEBUG {
		Println(green(indentation + checkmark + msg))
	}
}

// IndentRed prints red with indentation
func IndentRed(msg string) {
	if len(msg) > 0 {
		lines := strings.Split(msg, "\n")

		if len(lines) == 1 {
			Println(red(indentation + indentation + msg))
			return
		}

		for _, line := range lines {
			Print(red(indentation + indentation + line + "\n"))
		}
	}
}

// Failure prints green with xmark preceding the message
func Failure(msg string) {
	if len(msg) > 0 {
		Println(red(indentation + xmark + msg))
	}
}

// Info prints cyan with bullet preceding the message
func Info(msg string) {
	// if check for debug
	if logLevel < INFO {
		return
	}

	if len(msg) > 0 {
		lines := strings.Split(msg, "\n")

		if len(lines) == 1 {
			Println(green(indentation + bullet + msg))
			return
		}

		Print("\n")
		for i, line := range lines {
			if i == 0 {
				Print(green(indentation + bullet + line + "\n"))
				continue
			}
			Print(green(indentation + indentation + line + "\n"))
		}
	}
}

// Debug prints debugging info if debug flag is set
func Debug(msg string) {
	// fmt.Println("-----------------------------------------------------------------")

	if len(msg) > 0 && logLevel == DEBUG {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			Print(indentation + indentation + line + "\n")
		}
	}
}

// DebugLn prints debugging info if debug flag is set
func DebugLn(msg string) {
	if logLevel >= DEBUG {
		fmt.Println()
	}
	Debug(msg)
}

// Println prints specified message with line before the message
func Println(msg string) {
	fmt.Println()
	fmt.Println(msg)
}

// Print prints specified message
func Print(msg string) {
	fmt.Print(msg)
}

// Log prints message
func Log(msg string) {
	if len(msg) > 0 {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			fmt.Print(indentation + indentation + line + "\n")
		}
	}
}

// Newline prints new line
func Newline() {
	fmt.Println()
}

func colorString(msg, c string) string {
	if _isTerminal {
		switch c {
		case "red":
			return color.RedString(msg)
		case "green":
			return color.GreenString(msg)
		case "magenta":
			return color.MagentaString(msg)
		case "cyan":
			return color.CyanString(msg)
		default:
			return msg
		}
	}
	return msg
}

func red(msg string) string {
	return colorString(msg, "red")
}

func green(msg string) string {
	return colorString(msg, "green")
}

func magenta(msg string) string {
	return colorString(msg, "magenta")
}

func cyan(msg string) string {
	return colorString(msg, "cyan")
}

// RecursiveName travels up the parent chain to retrieve name
func RecursiveName(cmd *cobra.Command) string {

	if cmd.HasParent() {
		return RecursiveName(cmd.Parent()) + " " + cmd.Name()
	}

	return cmd.Name()
}

// UpdateExampleOnChildren updates children's example using parent name recursively
func UpdateExampleOnChildren(currCmd *cobra.Command) {

	if currCmd.HasSubCommands() {
		for _, c := range currCmd.Commands() {
			UpdateExampleOnChildren(c)
		}
	}

	if currCmd.HasParent() {
		currCmd.Example = "  " + RecursiveName(currCmd) + " " + currCmd.Example
		return
	}

	if len(currCmd.Example) > 0 {
		currCmd.Example = "  " + currCmd.Name() + " " + currCmd.Example
	}
}

// ConvertKeyValuePairArgSliceToMap converts slice of key=value to map
func ConvertKeyValuePairArgSliceToMap(keyvals []string) map[string]string {
	data := map[string]string{}
	for _, keyval := range keyvals {
		s := strings.TrimSpace(keyval)
		tokens := strings.Split(s, "=")
		// avoid error
		if len(tokens) != 2 {
			continue
		}
		key := tokens[0]
		val := tokens[1]
		data[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}

	return data
}
