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
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

var logger *slog.Logger

func SetLogger(l *slog.Logger) {
	logger = l
}

type IpCmd struct {
	BaseCommand
}

func New(path string) *IpCmd {
	return &IpCmd{
		BaseCommand: BaseCommand{path: path},
	}
}

func (i *IpCmd) AddNetns(name string) error {
	return i.execute("netns", "add", name)
}

func (i *IpCmd) DelNetns(name string) error {
	pids, err := i.ListNetnsProcesses(name)
	if err != nil {
		return err
	}

	if len(pids) > 0 {
		var pidStrs []string
		for _, pid := range pids {
			pidStrs = append(pidStrs, strconv.Itoa(pid))
		}
		return fmt.Errorf("netns %s has running processes: %s", name, strings.Join(pidStrs, ", "))
	}

	return i.execute("netns", "del", name)
}

func (i *IpCmd) ListNetns() []string {
	data, _ := i.executeWithOutput("netns", "list")

	var netns []string
	for _, line := range strings.Split(data, "\n") {
		name := strings.Split(line, " ")[0]
		netns = append(netns, name)
	}

	return netns
}

func (i *IpCmd) SetNetns(name string, netns string) error {
	return i.execute("link", "set", name, "netns", netns)
}

func (i *IpCmd) ListNetnsProcesses(netns string) ([]int, error) {
	out, err := i.executeWithOutput("netns", "pids", netns)
	if err != nil {
		return nil, err
	}

	var pids []int
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}

		pid, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}

func (i *IpCmd) ExistsNetns(name string) bool {
	return slices.Contains(i.ListNetns(), name)
}

func (i *IpCmd) IntoNetns(netns string) *IpCmdWithNetns {
	if netns == "" {
		return nil
	}

	ip := IpCmdWithNetns{
		netns: netns,
		BaseCommand: BaseCommand{
			path:        i.path,
			prependArgs: []string{"netns", "exec", netns, i.path},
		},
	}
	return &ip
}

type IpCmdWithNetns struct {
	netns string
	BaseCommand
}

func (i *IpCmdWithNetns) InNetns() bool {
	return true
}

func (i *IpCmdWithNetns) Netns() string {
	return i.netns
}

func (i *IpCmdWithNetns) ExecuteCommand(args ...string) (string, error) {
	return i.executeWithOutput(args...)
}
