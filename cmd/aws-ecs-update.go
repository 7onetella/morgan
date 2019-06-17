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
	"strings"

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsUpdateCmdCluster string
var ecsUpdateCmdDesiredCount int64
var ecsUpdateCmdTimeout int64
var ecsUpdateCmdWaitForServiceStable bool

var ecsUpdateCmd = &cobra.Command{
	Use:     "update-service <service name> <docker tags>",
	Short:   "Updates ecs",
	Long:    `Updates ecs`,
	Example: "foo-svc 1.0.0 --cluster api-cluster",
	Aliases: []string{"update"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsUpdateCmdCluster
		service := args[0]
		taskdef := ""
		tags := args[1:]

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

		// if --desired-count is not specified use the current count from service
		if ecsUpdateCmdDesiredCount == 0 {
			ecsUpdateCmdDesiredCount = *result.Services[0].DesiredCount
		}

		result2, err := ecsw.DescribeTaskDefinition(taskdef)
		ExitOnError(err, "describing task definition")

		for i, cd := range result2.TaskDefinition.ContainerDefinitions {
			imageCurr := *cd.Image
			colonIndex := strings.LastIndex(imageCurr, ":")
			imageBase := imageCurr[:colonIndex]
			// update the image with new docker tag
			imageCurr = imageBase + tags[i]
		}

		result3, err := ecsw.RegisterTaskDefinition(result2.TaskDefinition)
		ExitOnError(err, "registering task definition")

		_, err = ecsw.UpdateService(cluster, service, *result3.TaskDefinition.TaskDefinitionArn, ecsUpdateCmdDesiredCount)
		ExitOnError(err, "updating service")

		if ecsUpdateCmdWaitForServiceStable {
			err = ecsw.ServiceStable(cluster, service, ecsUpdateCmdTimeout)
			ExitOnError(err, "service stable")
		}

		Success("updating service")

	},
}

func init() {

	ecsCmd.AddCommand(ecsUpdateCmd)

	flags := ecsUpdateCmd.Flags()

	flags.StringVar(&ecsUpdateCmdCluster, "cluster", "", "ecs cluster")

	flags.Int64Var(&ecsUpdateCmdDesiredCount, "desired-count", 1, "desired count")

	flags.Int64Var(&ecsUpdateCmdTimeout, "timeout", 300, "timeout")

	flags.BoolVarP(&ecsUpdateCmdWaitForServiceStable, "service-stable", "w", false, "waits for service to become stable")

}
