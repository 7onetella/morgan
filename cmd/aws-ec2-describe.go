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
	"strings"

	"github.com/7onetella/morgan/tools/awsapi/ec2w"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ec2DescribeCmdName string

var ec2DescribeCmd = &cobra.Command{
	Use:     "describe",
	Short:   "Describes ec2 instances",
	Long:    `Describes ec2 instances`,
	Example: "",
	Aliases: []string{"describe-instances"},
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		Newline()

		resp, err := ec2w.DescribeInstanceByTagAndValue("", "")
		ExitOn(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "State", "Priv IP", "Instance ID"})

		for _, r := range resp.Reservations {
			for _, i := range r.Instances {
				var name string
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						name = *t.Value
					}
				}
				if len(ec2DescribeCmdName) > 0 && !strings.Contains(name, ec2DescribeCmdName) {
					continue
				}
				table.Append([]string{name, string(i.State.Name), *i.PrivateIpAddress, *i.InstanceId})
			}
		}
		table.Render()
	},
}

func init() {

	ec2Cmd.AddCommand(ec2DescribeCmd)

	flags := ec2DescribeCmd.Flags()

	flags.StringVar(&ec2DescribeCmdName, "name", "", "optional: ec2 name")

}
