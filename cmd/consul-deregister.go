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
	"github.com/spf13/cobra"
)

var consulDeregisterCmdClientAddr string

var consulDeregisterCmd = &cobra.Command{
	Use:     "deregister <name> <addr>",
	Short:   "Deregisters a Consul service instance",
	Example: "deregister api-service 10.10.10.10:34299",
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {

		CheckArgs(args, 2)

		name := args[0]
		addr := args[1]

		serviceID := name + "-" + addr
		clientaddr := consulDeregisterCmdClientAddr

		consulapi.Deregister(serviceID, clientaddr)

	},
}

func init() {

	consulCmd.AddCommand(consulDeregisterCmd)

	flags := consulDeregisterCmd.Flags()

	flags.StringVar(&consulDeregisterCmdClientAddr, "client", "169.254.1.1:8500", "consul client addr")

}
