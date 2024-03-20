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
	"netnsplan/config"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply netns networks configuration to running system",
	Long:  "Apply netns networks configuration to running system",
	RunE: func(cmd *cobra.Command, args []string) error {
		for netns, values := range cfg.Netns {
			if ip.ExistsNetns(netns) {
				slog.Warn("netns is already exists", "name", netns)
			} else {
				slog.Info("create netns", "name", netns)
				err := ip.AddNetns(netns)
				if err != nil {
					return err
				}
			}

			err := SetupLoopback(netns)
			if err != nil {
				return err
			}

			err = SetupEthernets(netns, values.Ethernets)
			if err != nil {
				return err
			}

			err = SetupDummyDevices(netns, values.DummyDevices)
			if err != nil {
				return err
			}

			err = SetupVethDevices(netns, values.VethDevices)
			if err != nil {
				return err
			}

			err = RunPostScript(netns, values.PostScript)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

type IpCommand interface {
	SetLinkUp(name string) error
	AddAddress(name, address string) error
	AddRoute(name, to, via string) error
	InNetns() bool
	Netns() string
}

func SetupDevice(ip IpCommand, name string, addresses []string, routes []config.Route) error {
	err := SetLinkUp(ip, name)
	if err != nil {
		return err
	}

	if ip.InNetns() {
		slog.Info("add addresses", "name", name, "addresses", addresses, "netns", ip.Netns())
	} else {
		slog.Info("add addresses", "name", name, "addresses", addresses)
	}
	for _, address := range addresses {
		err := ip.AddAddress(name, address)
		if err != nil {
			return err
		}
	}

	for _, route := range routes {
		if ip.InNetns() {
			slog.Info("add route", "name", name, "to", route.To, "via", route.Via, "netns", ip.Netns())
		} else {
			slog.Info("add route", "name", name, "to", route.To, "via", route.Via)
		}
		err := ip.AddRoute(name, route.To, route.Via)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetLinkUp(ip IpCommand, name string) error {
	if ip.InNetns() {
		slog.Info("link up", "name", name, "netns", ip.Netns())
	} else {
		slog.Info("link up", "name", name)
	}

	return ip.SetLinkUp(name)
}

func SetupLoopback(netns string) error {
	return SetLinkUp(ip.IntoNetns(netns), "lo")
}

func SetNetns(name string, netns string) error {
	slog.Info("set netns", "name", name, "netns", netns)
	return ip.SetNetns(name, netns)
}

func SetupEthernets(netns string, ethernets map[string]config.Ethernet) error {
	for name, values := range ethernets {
		err := SetNetns(name, netns)
		if err != nil {
			return err
		}

		err = SetupDevice(ip.IntoNetns(netns), name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupDummyDevices(netns string, devices map[string]config.Ethernet) error {
	for name, values := range devices {
		n := ip.IntoNetns(netns)

		slog.Info("add dummy device", "name", name, "netns", netns)
		err := n.AddDummyDevice(name)
		if err != nil {
			return err
		}

		err = SetupDevice(n, name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupVethDevices(netns string, devices map[string]config.VethDevice) error {
	for name, values := range devices {
		peerName := values.Peer.Name
		peerNetns := values.Peer.Netns

		slog.Info("add veth device", "name", name, "peer name", peerName)
		err := ip.AddVethDevice(name, peerName)
		if err != nil {
			return err
		}

		err = SetNetns(name, netns)
		if err != nil {
			return err
		}

		n := ip.IntoNetns(netns)
		err = SetupDevice(n, name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}

		if peerNetns != "" {
			err = SetNetns(peerName, peerNetns)
			if err != nil {
				return err
			}
			n := ip.IntoNetns(peerNetns)
			err = SetupDevice(n, peerName, values.Peer.Addresses, values.Peer.Routes)
			if err != nil {
				return err
			}
		} else {
			err = SetupDevice(ip, peerName, values.Peer.Addresses, values.Peer.Routes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func RunPostScript(netns string, script string) error {
	if script == "" {
		return nil
	}

	n := ip.IntoNetns(netns)
	slog.Info("run post script", "netns", netns, "script", script)
	out, err := n.ExecuteCommand(script)
	if err != nil {
		return err
	}
	slog.Debug("post script output", "netns", netns, "script", script, "output", out)

	return nil
}
