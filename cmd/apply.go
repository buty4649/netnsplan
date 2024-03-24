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
	"netnsplan/config"
	"netnsplan/iproute2"
	"slices"

	"github.com/spf13/cobra"
)

var alwaysRunPostScript bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply netns networks configuration to running system",
	Long:  "Apply netns networks configuration to running system",
	RunE: func(cmd *cobra.Command, args []string) error {
		for netns, values := range cfg.Netns {
			needPostScript := true
			if ip.ExistsNetns(netns) {
				slog.Warn("netns is already exists", "name", netns)
				needPostScript = false
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

			if needPostScript || alwaysRunPostScript {
				err = RunPostScript(netns, values.PostScript)
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().BoolVarP(&alwaysRunPostScript, "always-run-post-script", "R", false, "always run post-script. by default, runs only when a netns is created.")
}

type IpCommand interface {
	SetLinkUp(name string) error
	ShowLink(name string) (*iproute2.Link, error)
	ShowInterface(name string) (*iproute2.InterfaceInfo, error)
	AddAddress(name, address string) error
	ShowRoutes(name string) (iproute2.Routes, error)
	AddRoute(name, to, via string) error
	InNetns() bool
	Netns() string
}

func SetupDevice(ip IpCommand, name string, addresses []string, routes []config.Route) error {
	err := SetLinkUp(ip, name)
	if err != nil {
		return err
	}

	iface, err := ip.ShowInterface(name)
	if err != nil {
		return err
	}

	var intAddrs []string
	for _, i := range iface.AddrInfo {
		intAddrs = append(intAddrs, fmt.Sprintf("%s/%d", i.Local, i.Prefixlen))
	}

	for _, address := range addresses {
		if slices.Contains(intAddrs, address) {
			slog.Debug("address is already exists", "name", name, "address", address)
			continue
		}

		if ip.InNetns() {
			slog.Info("add addresses", "name", name, "address", address, "netns", ip.Netns())
		} else {
			slog.Info("add addresses", "name", name, "address", address)
		}

		err = ip.AddAddress(name, address)
		if err != nil {
			return err
		}
	}

	rt, err := ip.ShowRoutes(name)
	if err != nil {
		return err
	}

	for _, route := range routes {
		slog.Debug("route", "name", name, "to", route.To, "via", route.Via, "rt", rt)
		if slices.ContainsFunc(rt, func(r iproute2.Route) bool {
			return r.Dst == route.To && r.Gateway == route.Via
		}) {
			slog.Debug("route is already exists", "name", name, "to", route.To, "via", route.Via)
			continue
		}

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
	link, err := ip.ShowLink(name)
	if err != nil {
		return err
	}

	if slices.Contains(link.Flags, "UP") {
		slog.Debug("link is already up", "name", name)
		return nil
	}

	if ip.InNetns() {
		slog.Info("link up", "name", name, "netns", ip.Netns())
	} else {
		slog.Info("link up", "name", name)
	}

	return ip.SetLinkUp(name)
}

func SetNetns(name string, netns string) error {
	slog.Info("set netns", "name", name, "netns", netns)
	return ip.SetNetns(name, netns)
}

func SetupLoopback(netns string) error {
	return SetLinkUp(ip.IntoNetns(netns), "lo")
}

func SetupEthernets(netns string, ethernets map[string]config.Ethernet) error {
	n := ip.IntoNetns(netns)
	for name, values := range ethernets {
		_, err := n.ShowLink(name)
		if err == nil {
			slog.Debug("device is already exists in netns", "name", name, "netns", netns)
		} else {
			if _, ok := err.(*iproute2.NotExistError); ok {
				err := SetNetns(name, netns)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		err = SetupDevice(n, name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupDummyDevices(netns string, devices map[string]config.Ethernet) error {
	n := ip.IntoNetns(netns)
	for name, values := range devices {
		_, err := n.ShowLink(name)
		if err == nil {
			slog.Debug("device is already exists in netns", "name", name, "netns", netns)
		} else {
			if _, ok := err.(*iproute2.NotExistError); !ok {
				return err
			} else {
				slog.Info("add dummy device", "name", name, "netns", netns)
				err := n.AddDummyDevice(name)
				if err != nil {
					return err
				}
			}
		}

		err = SetupDevice(n, name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupVethDevices(netns string, devices map[string]config.VethDevice) error {
	n := ip.IntoNetns(netns)
	for name, values := range devices {
		peerName := values.Peer.Name
		peerNetns := values.Peer.Netns

		// check if device is already exists in netns
		_, err := n.ShowLink(name)
		if err == nil {
			slog.Debug("device is already exists in netns", "name", name, "netns", netns)
		} else {
			if _, ok := err.(*iproute2.NotExistError); !ok {
				return err
			} else {
				// check if device is already exists in "default" netns
				_, e := ip.ShowLink(name)
				if e == nil {
					slog.Debug("device is already exists", "name", name)
				} else {
					if _, ok := err.(*iproute2.NotExistError); !ok {
						return err
					} else {
						slog.Info("add veth device", "name", name, "peer name", peerName)
						err := ip.AddVethDevice(name, peerName)
						if err != nil {
							return err
						}
					}
				}

				err = SetNetns(name, netns)
				if err != nil {
					return err
				}
			}
		}

		err = SetupDevice(n, name, values.Addresses, values.Routes)
		if err != nil {
			return err
		}

		if peerNetns != "" {
			n := ip.IntoNetns(peerNetns)

			_, err := n.ShowLink(peerName)
			if err == nil {
				slog.Debug("device is already exists in netns", "name", peerName, "netns", peerNetns)
			} else {
				if _, ok := err.(*iproute2.NotExistError); !ok {
					return err
				} else {
					err = SetNetns(peerName, peerNetns)
					if err != nil {
						return err
					}
				}
			}
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
