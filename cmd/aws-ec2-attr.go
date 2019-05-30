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
	"fmt"

	"github.com/7onetella/morgan/tools/awsapi/ec2w"
	"github.com/spf13/cobra"
)

// ec2StopCmd represents the start command
var ec2AttrCmd = &cobra.Command{
	Use:   "attr <instance name> <keyname|privateip>",
	Short: "Display ec2 instance attribute value",
	Long:  `Display ec2 instance attribute value`,
	Example: `$ attr nginx keyname
my_ec2_pem_keyname

$ attr nginx privateip
172.31.10.53
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := ec2w.DescribeInstanceByNameTag(args[0])
		ExitOn(err)

		attr := args[1]

		switch attr {
		case "keyname":
			fmt.Print(*resp.Reservations[0].Instances[0].KeyName)
		case "privateip":
			fmt.Print(*resp.Reservations[0].Instances[0].PrivateIpAddress)
		}

	},
}

func init() {
	ec2Cmd.AddCommand(ec2AttrCmd)
}

// We can create a bash script called aws-ssh with the content like the following to automate ssh connection call

/*
#!/bin/sh

name=$1

keyname=$(morgan aws ec2 attr $name keyname)
privateip=$(morgan aws ec2 attr $name privateip)

ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i ~/.aws/${keyname}.pem ubuntu@${privateip}

if [ "$?" != "0" ]; then
  ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i ~/.aws/${keyname}.pem ec2-user@${privateip}
fi
*/
