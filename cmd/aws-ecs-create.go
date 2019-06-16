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
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsCreateCmdCluster string
var ecsCreateCmdDesiredCount int64
var ecsCreateCmdTimeout int64

var ecsCreateCmd = &cobra.Command{
	Use:     "create <service name> <size> <port> <image>",
	Short:   "Creates ecs",
	Long:    `Creates ecs`,
	Example: "foo-svc small 8080 7onetella/ref-api:latest",
	Args:    cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsCreateCmdCluster
		service := args[0]
		// size := args[1]
		port, _ := strconv.Atoi(args[2])
		image := args[3]

		clusters := map[string]string{}

		if len(cluster) == 0 {
			var err error
			clusters, err = ecsw.GetClustersForService(service)
			ExitOnError(err, "getting clusters for service")
		}

		if len(clusters) > 0 {
			Failure("creating service")
			Info("service [" + service + "] already exists in the following cluster or clusters")
			Newline()
			for k := range clusters {
				Print(indentation + "  - " + k + "\n")
			}
			os.Exit(1)
		}

		taskdef := &ecs.TaskDefinition{
			ContainerDefinitions: []ecs.ContainerDefinition{
				ecs.ContainerDefinition{
					Name:              aws.String(service),
					Cpu:               aws.Int64(128),
					Memory:            aws.Int64(256),
					MemoryReservation: aws.Int64(256),
					PortMappings: []ecs.PortMapping{
						ecs.PortMapping{
							ContainerPort: aws.Int64(int64(port)),
							HostPort:      aws.Int64(0),
							Protocol:      "tcp",
						},
					},
					Image: aws.String(image),
				},
			},
			Family: aws.String(service),
		}

		result, err := ecsw.RegisterTaskDefinition(taskdef)
		ExitOnError(err, "registering task definition")

		_, err = ecsw.CreateService(cluster, service, *result.TaskDefinition.TaskDefinitionArn, ecsCreateCmdDesiredCount)
		ExitOnError(err, "creating service")

		err = ecsw.ServiceStable(cluster, service, ecsCreateCmdTimeout)
		ExitOnError(err, "service stable")

		Success("creating service")

	},
}

func init() {

	ecsCmd.AddCommand(ecsCreateCmd)

	flags := ecsCreateCmd.Flags()

	flags.StringVar(&ecsCreateCmdCluster, "cluster", "", "ecs cluster")

	flags.Int64Var(&ecsCreateCmdDesiredCount, "desired-count", 1, "desired count")

	flags.Int64Var(&ecsCreateCmdTimeout, "timeout", 300, "timeout")

}
