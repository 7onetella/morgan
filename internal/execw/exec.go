package execw

// MIT License

// Copyright (c) 2019 7onetella

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"bufio"
	"fmt"
	"os/exec"
)

// Execute execute
func Execute(args []string) (string, error) {
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(output), err
}

func ExecuteRetrunStdoutErrOut(args []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(output), err
}

// Exec exec.Command
func Exec(args []string) (string, string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		fmt.Println("terminated early: " + err.Error())
		return "", "", err
	}

	var stdout string
	var errout string

	go func() {
		scanner := bufio.NewScanner(stdOut)
		for scanner.Scan() {
			// fmt.Println(indentation + scanner.Text())
			stdout = stdout + scanner.Text() + "\n"
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdErr)
		for scanner.Scan() {
			// fmt.Println(indentation + scanner.Text())
			errout = errout + scanner.Text() + "\n"
		}
	}()

	// if exec.Command calls a service that blocks then this code will never be reached
	err = cmd.Wait()

	return stdout, errout, err
}
