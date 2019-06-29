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

var ec2StartCmdTag string

// ec2StartCmd represents the start command
var ec2StartCmd = &cobra.Command{
	Use:   "start <instance names>",
	Short: "Starts ec2",
	Long: `Starts ec2

	you can start ec2 instances by "Name" tag:

	morgan aws ec2 start web-server api-server db-server

	or you can specify --tag flag to specify tag name and tag value

	morgan aws ec2 start --tag group:web-stack

`,
	Example: "start nginx redis",
	Aliases: []string{"start-instances"},
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		Newline()

		// retrieve ec2 instances by Name tag
		instanceIDs, err := ec2w.GetInstanceIDsByNames(args)
		ExitOnError(err, `retrieving by tag "Name"`)

		instanceIDSlice := []string{}
		for k := range instanceIDs {
			instanceIDSlice = append(instanceIDSlice, k)
		}

		// retrieve ec2 instances by tagName:tagValue if specified
		if len(ec2StartCmdTag) > 0 {

			tokens := strings.Split(ec2StartCmdTag, ":")
			if len(tokens) == 2 {
				tagName := strings.TrimSpace(tokens[0])
				tagValue := strings.TrimSpace(tokens[1])

				result, err := ec2w.DescribeInstanceByTagAndValue(tagName, tagValue)
				ExitOnError(err, "retrieving by tag name and tag value")

				for _, r := range result.Reservations {
					for _, i := range r.Instances {
						instanceIDSlice = append(instanceIDSlice, *i.InstanceId)

						for _, t := range i.Tags {
							if *t.Key == "Name" {
								instanceIDs[*i.InstanceId] = *t.Value
							}
						}
					}
				}
			}
		}

		// is there anything to start
		if len(instanceIDSlice) == 0 {
			Failure("there is no ec2 instances to start. check your parameters.")
			return
		}

		resp, err := ec2w.StartInstances(instanceIDSlice)
		ExitOn(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Instance ID", "Prev", "Current"})

		for _, c := range resp.StartingInstances {
			table.Append([]string{instanceIDs[*c.InstanceId], *c.InstanceId, string(c.PreviousState.Name), string(c.CurrentState.Name)})
		}
		table.Render()
	},
}

func init() {

	ec2Cmd.AddCommand(ec2StartCmd)

	flags := ec2StartCmd.Flags()

	flags.StringVar(&ec2StartCmdTag, "tag", "", "optional: ec2 tag")

}
