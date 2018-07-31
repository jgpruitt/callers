// MIT License
//
// Copyright (c) 2018 John Pruitt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package callers provides a simple and easy way to get a stack trace
package callers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Frame contains file/function/line information for a stack frame
type Frame struct {
	// the source file where the call was made
	File string
	// the line in the source file where the call was made
	Line int
	// the function where the call was made
	Function string
}

// String returns a simple string representation of the Frame
func (f *Frame) String() string {
	return fmt.Sprintf("File: %s Line: %d Function: %s", f.File, f.Line, f.Function)
}

// String prints a slice of *Frames to a string one per line.
// "indent" is printed at the beginning of each line
func String(trace []*Frame, indent string) string {
	var buf = &bytes.Buffer{}
	for _, frame := range trace {
		fmt.Fprint(buf, indent)
		fmt.Fprintln(buf, frame)
	}
	return buf.String()
}

// Callers is a handy wrapper around runtime.Callers and
// runtime.CallersFrames. The argument "skip" is the number
// of stack frames to skip before collecting frames. The
// "depth" argument is the number of stack frames to collect.
func Callers(skip, depth int) (trace []*Frame) {
	if skip < 0 {
		skip = 0
	}
	if depth <= 0 {
		depth = 10
	}

	trace = make([]*Frame, 0)
	var pc = make([]uintptr, depth)
	var n = runtime.Callers(skip, pc)
	var fs = runtime.CallersFrames(pc[:n])
	var f, ok = fs.Next()
	for ok {
		var frame = &Frame{
			Line:     f.Line,
			Function: f.Function,
		}
		var file = filepath.ToSlash(f.File)
		if n := strings.LastIndex(file, "/src/"); n > 0 {
			file = file[n+5:]
		} else {
			file = filepath.Base(file)
		}
		frame.File = file

		trace = append(trace, frame)
		f, ok = fs.Next()
	}
	return
}
