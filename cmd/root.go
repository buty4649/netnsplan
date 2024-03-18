/*
Copyright © 2024 buty4649

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"log/slog"
	"os"

	"netnsplan/config"
	"netnsplan/iproute2"
	"netnsplan/version"

	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
)

var cfgFilePath string
var ipCmdPath string
var debug bool

var cfg *config.Config
var ip *iproute2.Iproute2

var rootCmd = &cobra.Command{
	Use:     "netnsplan",
	Version: version.Version,
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logLevel := slog.LevelInfo
		if debug {
			logLevel = slog.LevelDebug
		}
		slog.SetDefault(slog.New(slogor.NewHandler(os.Stdout, slogor.Options{
			TimeFormat: "",
			Level:      logLevel,
			ShowSource: false,
		})))

		var err error
		cfg, err = config.LoadConfig(cfgFilePath)
		if err != nil {
			return err
		}

		ip = iproute2.New(ipCmdPath)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", "./netnsplan.yaml", "config file")
	rootCmd.PersistentFlags().StringVar(&ipCmdPath, "cmd", "/bin/ip", "ip command path")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
}
