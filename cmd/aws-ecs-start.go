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
	"errors"

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsStartCmdCluster string
var ecsStartCmdTimeout int64
var ecsStartCmdDesiredCount int64

var ecsStartCmd = &cobra.Command{
	Use:     "start <service name>",
	Short:   "Starts ecs",
	Long:    `Starts ecs`,
	Example: "foo-svc -c api-cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsStartCmdCluster
		service := args[0]
		taskdef := ""

		// if cluster is not specified, then assume there is only one cluster and use that cluster
		if len(cluster) == 0 {
			clusters := GetClustersForService(service)

			CheckForClusterAmbiguity(clusters)

			cluster = GetClusterForService(clusters, service)
		}

		result, err := ecsw.DescribeServices(cluster, service)
		ExitOnError(err, "describing services")
		if len(result.Services) == 0 {
			ExitOnError(errors.New("search result count 0"), "finding service")
		}
		taskdef = *result.Services[0].TaskDefinition

		_, err = ecsw.UpdateService(cluster, service, taskdef, ecsStartCmdDesiredCount)
		ExitOnError(err, "updating service with specified desired count")

		err = ecsw.ServiceStable(cluster, service, ecsStartCmdTimeout)
		ExitOnError(err, "service stable")

		Success("starting service")

	},
}

func init() {

	ecsCmd.AddCommand(ecsStartCmd)

	flags := ecsStartCmd.Flags()

	flags.StringVarP(&ecsStartCmdCluster, "cluster", "c", "", "ecs cluster")

	flags.Int64VarP(&ecsStartCmdTimeout, "timeout", "t", 300, "timeout")

	flags.Int64Var(&ecsStartCmdDesiredCount, "desired-count", 1, "desired count")

}
