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
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"syscall"
)

type BaseCommand struct {
	path    string
	prepend []string
}

type CommandOut struct {
	Stdout string
	Stderr string
}

func (b *BaseCommand) run(args ...string) error {
	_, err := b.runIpCommand(args...)
	return err
}

func (b *BaseCommand) runIpCommand(args ...string) (string, error) {
	cmd := append([]string{b.path}, args...)
	out, err := b.runCommand(cmd, nil)
	if err == nil {
		if out.Stderr != "" {
			slog.Warn("ip command warning", "msg", out.Stderr)
		}

		return out.Stdout, nil
	}

	msg := err.Error()
	if strings.Contains(msg, "Operation not permitted") {
		return "", &OperationNotPermittedError{Msg: msg}
	}

	if strings.Contains(msg, "does not exist") {
		return "", &NotExistError{Msg: msg}
	}

	return "", &UnknownError{Msg: err.Error()}
}

func (b *BaseCommand) runCommand(cmd []string, input *string) (*CommandOut, error) {
	if b.prepend != nil {
		cmd = append(b.prepend, cmd...)
	}
	path := cmd[0]
	args := cmd[1:]

	if logger != nil {
		logger.Debug("exec", "cmd", path, "args", args)
	}

	c := exec.Command(path, args...)
	if input != nil {
		stdin, err := c.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, *input)
		}()
	}

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		exitErr, _ := err.(*exec.ExitError)
		status, _ := exitErr.Sys().(syscall.WaitStatus)
		return nil, &CommandError{
			ExitStatus: status.ExitStatus(),
			Msg:        stderr.String(),
		}
	}
	return &CommandOut{Stdout: stdout.String(), Stderr: stderr.String()}, nil
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

type OperState string

const (
	OperStateUp     OperState = "UP"
	OperStateDown   OperState = "DOWN"
	OperStateUnkwon OperState = "UNKNOWN"
)

type InterfaceInfo struct {
	Ifindex   int           `json:"ifindex"`
	Ifname    string        `json:"ifname"`
	Flags     []string      `json:"flags"`
	Mtu       int           `json:"mtu"`
	Qdisc     string        `json:"qdisc"`
	Operstate OperState     `json:"operstate"`
	Group     string        `json:"group"`
	Txqlen    int           `json:"txqlen"`
	LinkType  string        `json:"link_type"`
	Address   string        `json:"address"`
	Broadcast string        `json:"broadcast"`
	AddrInfo  []AddressInfo `json:"addr_info"`
}

type Interfaces []InterfaceInfo

func (b *BaseCommand) ListInterfaces() (Interfaces, error) {
	data, err := b.runIpCommand("-json", "address", "show")
	if err != nil {
		return nil, err
	}

	return unmarshalInterfacesData(data)
}

func (b *BaseCommand) ShowInterface(name string) (*InterfaceInfo, error) {
	data, err := b.runIpCommand("-json", "address", "show", "dev", name)
	if err != nil {
		return nil, err
	}

	i, err := unmarshalInterfacesData(data)
	if err != nil {
		return nil, err
	}

	return &i[0], err
}

func unmarshalInterfacesData(data string) (Interfaces, error) {
	var addresses Interfaces
	err := json.Unmarshal([]byte(data), &addresses)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

type Link struct {
	Ifindex   int       `json:"ifindex"`
	Ifname    string    `json:"ifname"`
	Flags     []string  `json:"flags"`
	Mtu       int       `json:"mtu"`
	Qdisc     string    `json:"qdisc"`
	Operstate OperState `json:"operstate"`
	Linkmode  string    `json:"linkmode"`
	Group     string    `json:"group"`
	Txqlen    int       `json:"txqlen"`
	LinkType  string    `json:"link_type"`
	Address   string    `json:"address"`
	Broadcast string    `json:"broadcast"`
}

type Links []Link

func (b *BaseCommand) ListLinks() (Links, error) {
	data, err := b.runIpCommand("-json", "link", "show")
	if err != nil {
		return nil, err
	}

	return unmarshalLinksData(data)
}
func (b *BaseCommand) ShowLink(name string) (*Link, error) {
	data, err := b.runIpCommand("-json", "link", "show", "dev", name)
	if err != nil {
		return nil, err
	}

	links, err := unmarshalLinksData(data)
	if err != nil {
		return nil, err
	}

	return &links[0], nil
}

func unmarshalLinksData(data string) (Links, error) {
	var links Links
	err := json.Unmarshal([]byte(data), &links)
	if err != nil {
		return nil, err
	}

	return links, nil
}

type Route struct {
	Dst      string   `json:"dst,omitempty"`
	Gateway  string   `json:"gateway,omitempty"`
	Dev      string   `json:"dev,omitempty"`
	Type     string   `json:"type,omitempty"`
	Protocol string   `json:"protocol,omitempty"`
	Scope    string   `json:"scope,omitempty"`
	PrefSrc  string   `json:"prefsrc,omitempty"`
	Flags    []string `json:"flags,omitempty"`
}

type Routes []Route

func (b *BaseCommand) ListRoutes() (Routes, error) {
	data, err := b.runIpCommand("-json", "route", "show")
	if err != nil {
		return nil, err
	}

	return unmarshalRoutesData(data)
}

func (b *BaseCommand) ShowRoutes(name string) (Routes, error) {
	data, err := b.runIpCommand("-json", "route", "show", "dev", name)
	if err != nil {
		return nil, err
	}

	return unmarshalRoutesData(data)
}

func unmarshalRoutesData(data string) (Routes, error) {
	var routes []Route
	err := json.Unmarshal([]byte(data), &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
