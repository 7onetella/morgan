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

	"github.com/7onetella/morgan/tools/awsapi/route53"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/7onetella/morgan/tools/awsapi/ec2w"
	"github.com/spf13/cobra"
)

// dnsUpdateCmd represents the dns update command
var dnsUpdateCmd = &cobra.Command{
	Use:     "update <instance name> <A record name>",
	Short:   "Updates dns record",
	Long:    `Display dns record`,
	Example: `$`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		instanceName := args[0]
		dnsARecord := args[1]

		resp, err := ec2w.DescribeInstanceByNameTag(instanceName)
		ExitOn(err)

		if len(resp.Reservations) == 0 {
			Failure("looking up instance by tag name = " + instanceName)
			os.Exit(1)
		}

		privateIP := *resp.Reservations[0].Instances[0].PrivateIpAddress

		zresp, err := route53.ListHostedZones()
		ExitOn(err)

		terms := strings.Split(dnsARecord, ".")
		tSize := len(terms)
		rootDomain := terms[tSize-2] + "." + terms[tSize-1] + "."

		var hostedZoneID string
		for _, zone := range zresp.HostedZones {
			if *zone.Name == rootDomain {
				s := aws.StringValue(zone.Id)
				i := strings.LastIndex(s, "/")
				hostedZoneID = s[i+1:]
				break
			}
		}

		_, err = route53.ARecordUpsert(dnsARecord, privateIP, hostedZoneID)
		ExitOnError(err, "Updating A Record")

	},
}

func init() {
	dnsCmd.AddCommand(dnsUpdateCmd)
}
