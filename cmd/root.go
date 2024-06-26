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
	"fmt"
	"log/slog"
	"os"

	"netnsplan/config"
	"netnsplan/iproute2"
	"netnsplan/version"

	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
)

type Flags struct {
	ConfigDir    string
	IpCmdPath    string
	Debug, Quiet bool
}

var flags Flags

var cfg *config.Config
var ip *iproute2.IpCmd

var rootCmd = &cobra.Command{
	Use:          "netnsplan",
	Short:        "Easily automate Linux netns networks and configurations via YAML",
	Long:         "Easily automate Linux netns networks and configurations via YAML",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.ValidateFlagGroups()
		if err != nil {
			return err
		}

		var logLevel slog.Level
		if flags.Debug {
			logLevel = slog.LevelDebug
		}
		if flags.Quiet {
			logLevel = slog.LevelError
		}
		logger := slog.New(slogor.NewHandler(os.Stdout, slogor.Options{
			TimeFormat: "",
			Level:      logLevel,
			ShowSource: false,
		}))
		slog.SetDefault(logger)

		cfg, err = config.LoadYamlFiles(flags.ConfigDir)
		if err != nil {
			return err
		}

		ip = iproute2.New(flags.IpCmdPath)
		iproute2.SetLogger(logger)
		return nil
	},
}

func Execute() {
	if version.Version != "dev" {
		rootCmd.Version = fmt.Sprintf("v%s", version.Version)
	} else {
		rootCmd.Version = version.Version
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.ConfigDir, "config-dir", "d", "/etc/netnsplan", "config file directory")
	rootCmd.PersistentFlags().StringVar(&flags.IpCmdPath, "cmd", "/bin/ip", "ip command path")

	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "debug mode")
	rootCmd.MarkFlagsMutuallyExclusive("debug", "quiet")
}
