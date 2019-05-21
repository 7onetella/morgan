package cmd

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
	"github.com/7onetella/morgan/tools/consulapi"
	"github.com/7onetella/morgan/tools/ecsapi"
	"github.com/spf13/cobra"
)

var consulRegisterCmdClientAddr string

var consulRegisterCmd = &cobra.Command{
	Use:     "register <name> <urlprefix> <check>",
	Short:   "Registers a service to Consul",
	Example: "register api-service www.example.com/api /api/health",
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {

		CheckArgs(args, 3)

		name := args[0]
		if len(name) == 0 {
			var err error
			name, err = ecsapi.GetContainerName()
			ExitOnError(err, "container name")
		}

		host, err := ecsapi.GetHost()
		ExitOnError(err, "retrieving host")

		clientaddr := consulRegisterCmdClientAddr

		port, _ := ecsapi.GetHostPort()
		ExitOnError(err, "retrieving port")

		check := args[2]

		tags := []string{"urlprefix-" + args[1]}

		consulapi.Register(name, host, port, check, tags, clientaddr)

	},
}

func init() {

	consulCmd.AddCommand(consulRegisterCmd)

	flags := consulRegisterCmd.Flags()

	flags.StringVar(&consulRegisterCmdClientAddr, "client", "169.254.1.1:8500", "consul client addr")

}
