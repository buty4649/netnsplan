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

import "testing"

func TestNotExistError(t *testing.T) {
	msg := "Device \"test\" does not exist"
	err := &NotExistError{Msg: msg}

	if got := err.Error(); got != msg {
		t.Errorf("NotExistError.Error() = %v, want %v", got, msg)
	}
}

func TestOperationNotPermittedError(t *testing.T) {
	msg := "mount --make-shared /run/testnetns failed: Operation not permitted"
	err := &OperationNotPermittedError{Msg: msg}

	if got := err.Error(); got != msg {
		t.Errorf("OperationNotPermittedError.Error() = %v, want %v", got, msg)
	}
}

func TestCommandError(t *testing.T) {
	exitStatus := 1
	msg := "Command failed\n"
	want := "Command failed (exit status: 1)"
	err := &CommandError{ExitStatus: exitStatus, Msg: msg}

	if got := err.Error(); got != want {
		t.Errorf("CommandError.Error() = %v, want %v", got, want)
	}

	exitStatus = 2
	msg = "Error: Nexthop has invalid gateway.\n"
	want = "Nexthop has invalid gateway. (exit status: 2)"
	err = &CommandError{ExitStatus: exitStatus, Msg: msg}

	if got := err.Error(); got != want {
		t.Errorf("CommandError.Error() = %v, want %v", got, want)
	}
}
