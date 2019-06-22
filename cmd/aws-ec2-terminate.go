// Copyright © 2019 Seven OneTella<7onetella@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"

	"github.com/7onetella/morgan/tools/awsapi/ec2w"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ec2TerminateCmd = &cobra.Command{
	Use:     "terminate <instance names>",
	Short:   "Terminates ec2",
	Long:    `Terminates ec2`,
	Example: "nginx redis",
	Aliases: []string{"terminate-instances"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		Newline()

		instanceIDs, err := ec2w.GetInstanceIDsByNames(args)
		ExitOn(err)

		instanceIDSlice := []string{}
		for k := range instanceIDs {
			instanceIDSlice = append(instanceIDSlice, k)
		}

		resp, err := ec2w.TerminateInstances(instanceIDSlice)
		ExitOn(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Instance ID", "Prev", "Current"})

		for _, c := range resp.TerminatingInstances {
			table.Append([]string{instanceIDs[*c.InstanceId], *c.InstanceId, string(c.PreviousState.Name), string(c.CurrentState.Name)})
		}
		table.Render()
	},
}

func init() {
	ec2Cmd.AddCommand(ec2TerminateCmd)
}
