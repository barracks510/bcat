// Copyright © 2016 Dennis Chen <barracks510@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barracks510/bcat/bcatlib"
)

var (
	cfgFile string

	options bcatlib.Options
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bcat",
	Short: "pipe to browser utility",
	Long: `Usage: bcat [-htp] [-a] [-b <browser>] [-T <title>] [<file>]...
       bcat [-htp] [-a] [-b <browser>] [-T <title>] -c command...
       btee <options> [<file>]...
Pipe to browser utility. Read standard input, possibly one or more <file>s,
and write concatenated / formatted output to browser. When invoked as btee,
also write all input back to standard output.

Display options:
  -b, --browser=<browser>    open <browser> instead of system default browser
  -T, --title=<text>         use <text> as the browser title
  -a, --ansi                 convert ANSI (color) escape sequences to HTML

Input format (auto detected by default):
  --html                 input is already HTML encoded, doc or fragment
  --text                 input is unencoded text

Misc options:
  -c, --command              read the standard output of command
  -p, --persist              serve until interrupted, allowing reload
  -d, --debug                enable verbose debug logging on stderr`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var buffer []byte
		if len(args) == 0 || args[0] == "-" {
			buffer, _ = ioutil.ReadAll(os.Stdin)
		}
		browserCommand := viper.GetString("BCAT_COMMAND")
		b, err := bcatlib.NewBrowser(options.Browser, browserCommand)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		s, err := bcatlib.NewServer(func(w http.ResponseWriter, r *http.Request) {
			w.Write(buffer)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		if err := b.Open(s.Url()); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", s.Serve())
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.SetHelpTemplate(RootCmd.Long)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bcat.yaml)")

	RootCmd.PersistentFlags().BoolVar(&options.Html, "html", false, "")
	RootCmd.PersistentFlags().BoolVar(&options.Html, "text", false, "")
	RootCmd.PersistentFlags().StringVarP(&options.Browser, "browser", "b", "default", "")
	RootCmd.PersistentFlags().StringVarP(&options.Title, "title", "T", "", "")
	RootCmd.PersistentFlags().BoolVarP(&options.Ansi, "ansi", "a", false, "")
	RootCmd.PersistentFlags().BoolVarP(&options.Persist, "persist", "p", false, "")
	RootCmd.PersistentFlags().BoolVarP(&options.Command, "command", "c", false, "")
	RootCmd.PersistentFlags().BoolVarP(&options.Debug, "debug", "d", false, "")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".bcat") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
