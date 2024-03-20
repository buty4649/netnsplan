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
package iproute2

import (
	"encoding/json"
)

type ipCmdRunner interface {
	cmdPath() string
	runIpCommand(args ...string) (string, error)
}

type BaseCommand struct {
	runner ipCmdRunner
}

func (b *BaseCommand) run(args ...string) error {
	_, err := b.runner.runIpCommand(args...)
	return err
}

func (b *BaseCommand) AddLink(name string, linkType string, options ...string) error {
	args := append([]string{"link", "add", name, "type", linkType}, options...)
	return b.run(args...)
}

func (b *BaseCommand) DelLink(name string) error {
	return b.run("link", "del", name)
}

func (b *BaseCommand) AddDummyDevice(name string) error {
	return b.AddLink(name, "dummy")
}

func (b *BaseCommand) AddVethDevice(name string, peerName string) error {
	return b.AddLink(name, "veth", "peer", "name", peerName)
}

func (b *BaseCommand) SetLinkUp(name string) error {
	return b.run("link", "set", "dev", name, "up")
}

func (b *BaseCommand) AddAddress(name string, address string) error {
	return b.run("address", "add", address, "dev", name)
}

func (b *BaseCommand) DelAddress(name string, address string) error {
	return b.run("address", "del", address, "dev", name)
}

func (b *BaseCommand) AddRoute(name string, to string, via string) error {
	return b.run("route", "add", to, "via", via, "dev", name)
}

func (b *BaseCommand) DelRoute(name string, to string, via string) error {
	return b.run("route", "del", to, "via", via, "dev", name)
}

func (i *IpCmd) InNetns() bool {
	return false
}

func (i *IpCmd) Netns() string {
	return ""
}

type AddressInfo struct {
	Family            string `json:"family"`
	Local             string `json:"local"`
	Prefixlen         int    `json:"prefixlen"`
	Scope             string `json:"scope"`
	Label             string `json:"label"`
	ValidLifeTime     uint64 `json:"valid_life_time"`
	PreferredLifeTime uint64 `json:"preferred_life_time"`
}

type InterfaceInfo struct {
	Ifindex   int           `json:"ifindex"`
	Ifname    string        `json:"ifname"`
	Flags     []string      `json:"flags"`
	Mtu       int           `json:"mtu"`
	Qdisc     string        `json:"qdisc"`
	Operstate string        `json:"operstate"`
	Group     string        `json:"group"`
	Txqlen    int           `json:"txqlen"`
	LinkType  string        `json:"link_type"`
	Address   string        `json:"address"`
	Broadcast string        `json:"broadcast"`
	AddrInfo  []AddressInfo `json:"addr_info"`
}

type Addresses []InterfaceInfo

func (b *BaseCommand) ListAddresses() (*Addresses, error) {
	data, err := b.runner.runIpCommand("-json", "address", "show")
	if err != nil {
		return nil, err
	}

	var addresses Addresses
	err = json.Unmarshal([]byte(data), &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}

type Link struct {
	Ifindex   int      `json:"ifindex"`
	Ifname    string   `json:"ifname"`
	Flags     []string `json:"flags"`
	Mtu       int      `json:"mtu"`
	Qdisc     string   `json:"qdisc"`
	Operstate string   `json:"operstate"`
	Linkmode  string   `json:"linkmode"`
	Group     string   `json:"group"`
	Txqlen    int      `json:"txqlen"`
	LinkType  string   `json:"link_type"`
	Address   string   `json:"address"`
	Broadcast string   `json:"broadcast"`
}

type Links []Link

func (b *BaseCommand) ListLinks() (*Links, error) {
	data, err := b.runner.runIpCommand("-json", "link", "show")
	if err != nil {
		return nil, err
	}

	var links Links
	err = json.Unmarshal([]byte(data), &links)
	if err != nil {
		return nil, err
	}

	return &links, nil
}
