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
	"strings"
)

type NotExistError struct {
	Msg string
}

func (e *NotExistError) Error() string {
	return e.Msg
}

type OperationNotPermittedError struct {
	Msg string
}

func (e *OperationNotPermittedError) Error() string {
	return e.Msg
}

type UnmarshalError struct {
	Msg     string
	Content string
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprint(e.Msg, ": ", e.Content)
}

type UnknownError struct {
	Msg string
}

func (e *UnknownError) Error() string {
	return e.Msg
}

type CommandError struct {
	ExitStatus int
	Msg        string
}

func (e *CommandError) Error() string {
	msg := strings.TrimRight(e.Msg, "\n")
	msg = strings.TrimPrefix(msg, "Error: ")
	return fmt.Sprintf("%s (exit status: %d)", msg, e.ExitStatus)
}
