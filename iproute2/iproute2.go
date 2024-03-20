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
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

type Iproute2 struct {
	path     string
	netns    string
	useNetns bool
}

type Error struct {
	ExitStatus int
	Message    string
}

var logger *slog.Logger

func New(path string) *Iproute2 {
	return &Iproute2{
		path:     path,
		netns:    "",
		useNetns: false,
	}
}

func SetLogger(l *slog.Logger) {
	logger = l
}

func (i *Iproute2) AddLink(name string, linkType string, options ...string) error {
	args := []string{"link", "add", "name", name, "type", linkType}
	args = append(args, options...)
	return i.execute(args...)
}

func (i *Iproute2) DelLink(name string) error {
	return i.execute("link", "del", "name", name)
}

func (i *Iproute2) AddDummyDevice(name string) error {
	return i.AddLink(name, "dummy")
}

func (i *Iproute2) AddVethDevice(name string, peerName string) error {
	return i.AddLink(name, "veth", "peer", "name", peerName)
}

func (i *Iproute2) SetLinkUp(name string) error {
	return i.execute("link", "set", "dev", name, "up")
}

func (i *Iproute2) AddAddress(name string, address string) error {
	return i.execute("address", "add", address, "dev", name)
}

func (i *Iproute2) DelAddress(name string, address string) error {
	return i.execute("address", "del", address, "dev", name)
}

func (i *Iproute2) AddRoute(name string, to string, via string) error {
	return i.execute("route", "add", to, "via", via, "dev", name)
}

func (i *Iproute2) DelRoute(name string, to string, via string) error {
	return i.execute("route", "del", to, "via", via, "dev", name)
}

func (i *Iproute2) AddNetns(name string) error {
	return i.execute("netns", "add", name)
}

func (i *Iproute2) DelNetns(name string) error {
	pids, err := i.NetnsPid(name)
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

func (i *Iproute2) ListNetns() []string {
	data, _ := i.executeWithStdout("netns", "list")

	var netns []string
	for _, line := range strings.Split(data, "\n") {
		name := strings.Split(line, " ")[0]
		netns = append(netns, name)
	}

	return netns
}

func (i *Iproute2) SetNetns(name string, netns string) error {
	return i.execute("link", "set", name, "netns", netns)
}

func (i *Iproute2) IntoNetns(netns string, fn func() error) error {
	i.netns = netns
	i.useNetns = true

	err := fn()

	i.netns = ""
	i.useNetns = false
	return err
}

func (i *Iproute2) Netns() string {
	if i.useNetns {
		return i.netns
	}
	return ""
}

func (i *Iproute2) NetnsPid(netns string) ([]int, error) {
	out, err := i.executeWithStdout("netns", "pids", netns)
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

func (i *Iproute2) NetnsExists(name string) bool {
	return slices.Contains(i.ListNetns(), name)
}

func (e *Error) Error() string {
	msg := strings.TrimRight(e.Message, "\n")
	return fmt.Sprintf("%s (exit status: %d)", msg, e.ExitStatus)
}

func (i *Iproute2) execute(args ...string) error {
	_, err := i.executeWithStdout(args...)
	return err
}

func (i *Iproute2) executeWithStdout(args ...string) (string, error) {
	var cmdArgs []string

	if i.useNetns {
		cmdArgs = append(cmdArgs, "netns", "exec", i.netns)
		cmdArgs = append(cmdArgs, i.path)
	}
	cmdArgs = append(cmdArgs, args...)

	if logger != nil {
		logger.Debug("exec command", "path", i.path, "args", strings.Join(cmdArgs, " "))
	}

	cmd := exec.Command(i.path, cmdArgs...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()

	if err != nil {
		exitErr, _ := err.(*exec.ExitError)
		status, _ := exitErr.Sys().(syscall.WaitStatus)
		return "", &Error{
			ExitStatus: status.ExitStatus(),
			Message:    stderrBuf.String(),
		}
	}

	return stdoutBuf.String(), nil
}
