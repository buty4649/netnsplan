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
	"os/exec"
	"strings"
	"syscall"
)

type ipCmd struct {
	path string
}

func (i *ipCmd) runIpCommand(args ...string) (string, error) {
	return runCommand(i.path, args...)
}

type ipCmdWithNetns struct {
	path  string
	netns string
}

func (i *ipCmdWithNetns) runIpCommand(args ...string) (string, error) {
	cmd := append([]string{i.path}, args...)
	return i.runWithNetns(cmd...)
}

func (i *ipCmdWithNetns) runWithNetns(cmd ...string) (string, error) {
	cmdArgs := append([]string{"netns", "exec", i.netns}, cmd...)
	return runCommand(i.path, cmdArgs...)
}

type Error struct {
	ExitStatus int
	Message    string
}

func (e *Error) Error() string {
	msg := strings.TrimRight(e.Message, "\n")
	return fmt.Sprintf("%s (exit status: %d)", msg, e.ExitStatus)
}

func runCommand(path string, args ...string) (string, error) {
	if logger != nil {
		logger.Debug("exec", "cmd", path, "args", args)
	}

	cmd := exec.Command(path, args...)
	stdout, err := cmd.Output()
	if err != nil {
		exitErr, _ := err.(*exec.ExitError)
		status, _ := exitErr.Sys().(syscall.WaitStatus)
		stderr := string(exitErr.Stderr)
		return "", &Error{
			ExitStatus: status.ExitStatus(),
			Message:    stderr,
		}
	}
	return string(stdout), nil
}
