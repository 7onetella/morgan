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
	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsDeleteCmdCluster string

var ecsDeleteCmd = &cobra.Command{
	Use:     "delete-service <service name>",
	Short:   "Deletes ecs",
	Long:    `Deletes ecs`,
	Example: "foo-svc --cluster api-cluster",
	Aliases: []string{"delete"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsDeleteCmdCluster
		service := args[0]

		// if cluster is not specified, then assume there is only one cluster and use that cluster
		if len(cluster) == 0 {
			clusters := GetClustersForService(service)

			CheckForClusterAmbiguity(clusters)

			cluster = GetClusterForService(clusters, service)
		}

		result, err := ecsw.DescribeServices(cluster, service)
		ExitOnError(err, "describing services")
		taskdef := *result.Services[0].TaskDefinition

		result2, err := ecsw.DescribeTaskDefinition(taskdef)
		ExitOnError(err, "describing task definition")

		_, err = ecsw.UpdateService(cluster, service, *result2.TaskDefinition.TaskDefinitionArn, 0)
		ExitOnError(err, "updating service")

		_, err = ecsw.DeleteService(cluster, service)
		ExitOnError(err, "deleting services")

		Success("deleting service")

	},
}

func init() {

	ecsCmd.AddCommand(ecsDeleteCmd)

	flags := ecsDeleteCmd.Flags()

	flags.StringVarP(&ecsDeleteCmdCluster, "cluster", "c", "", "ecs cluster")

}
