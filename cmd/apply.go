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
	Use:     "apply",
	Short:   "Apply netns networks configuration to running system",
	Long:    "Apply netns networks configuration to running system",
	RunE: func(cmd *cobra.Command, args []string) error {
		for netns, values := range cfg.Netns {
			if ip.NetnsExists(netns) {
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
		}
		return nil
	},
}

func SetupDevice(name string, addresses []string) error {
	for _, address := range addresses {
		err := ip.AddAddress(name, address)
		if err != nil {
			return err
		}
	}

	err := ip.SetLinkUp(name)
	if err != nil {
		return err
	}

	return nil
}

func SetupLoopback(netns string) error {
	return ip.IntoNetns(netns, func() error {
		return ip.SetLinkUp("lo")
	})
}

func SetupEthernets(netns string, ethernets map[string]config.EthernetConfig) error {
	for name, values := range ethernets {
		err := ip.SetNetns(name, netns)
		if err != nil {
			return err
		}

		ip.IntoNetns(netns, func() error {
			return SetupDevice(name, values.Addresses)
		})
	}
	return nil
}

func SetupDummyDevices(netns string, devices map[string]config.EthernetConfig) error {
	for name, values := range devices {
		ip.IntoNetns(netns, func() error {
			err := ip.AddDummyDevice(name)
			if err != nil {
				return err
			}

			return SetupDevice(name, values.Addresses)
		})
	}
	return nil
}

func SetupVethDevices(netns string, devices map[string]config.VethDeviceConfig) error {
	for name, values := range devices {
		peerName := values.Peer.Name
		peerNetns := values.Peer.Netns

		err := ip.AddVethDevice(name, peerName)
		if err != nil {
			return err
		}

		err = ip.SetNetns(name, netns)
		if err != nil {
			return err
		}
		ip.IntoNetns(netns, func() error {
			return SetupDevice(name, values.Addresses)
		})

		if peerNetns != "" {
			err = ip.SetNetns(peerName, peerNetns)
			if err != nil {
				return err
			}
			ip.IntoNetns(netns, func() error {
				err = SetupDevice(peerName, values.Peer.Addresses)
				if err != nil {
					return err
				}
				return nil
			})
		} else {
			SetupDevice(peerName, values.Peer.Addresses)
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
