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
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsCreateCmdCluster string
var ecsCreateCmdService string
var ecsCreateCmdTaskDefinition string
var ecsCreateCmdDesiredCount int64
var ecsCreateCmdTimeout int64
var ecsCreateCmdWaitForServiceStable bool

var ecsCreateCmd = &cobra.Command{
	Use:     "create-service <service-name> <size> <port> <image>",
	Short:   "Creates ecs",
	Long:    `Creates ecs`,
	Example: "foo-svc small 8080 7onetella/ref-api:latest",
	Aliases: []string{"create"},
	Args:    cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsCreateCmdCluster
		service := args[0]
		if len(service) == 0 {
			service = ecsCreateCmdService
		}
		size := args[1]
		port, _ := strconv.Atoi(args[2])
		image := args[3]
		taskdef := ecsCreateCmdTaskDefinition

		// if cluster is not specified, then assume there is only one cluster and use that cluster
		if len(cluster) == 0 {
			clusters := GetClustersForService(service)

			CheckForClusterAmbiguity(clusters)

			cluster = GetClusterForService(clusters, service)
		}

		if len(taskdef) == 0 {
			cm := GetCPUAndMemory(size)
			taskdef = RegisterNewTaskDefinition(cm.CPU, cm.Memory, int64(port), service, image)
		}

		_, err := ecsw.CreateService(cluster, service, taskdef, ecsCreateCmdDesiredCount)
		ExitOnError(err, "creating service")

		if ecsCreateCmdWaitForServiceStable {
			err = ecsw.ServiceStable(cluster, service, ecsCreateCmdTimeout)
			ExitOnError(err, "service stable")
		}

		Success("creating service")

	},
}

func init() {

	ecsCmd.AddCommand(ecsCreateCmd)

	flags := ecsCreateCmd.Flags()

	flags.StringVarP(&ecsCreateCmdCluster, "cluster", "c", "", "ecs cluster")

	flags.StringVarP(&ecsCreateCmdService, "service", "s", "", "service name")

	flags.StringVar(&ecsCreateCmdTaskDefinition, "task-definition", "", "task definition")

	flags.Int64Var(&ecsCreateCmdDesiredCount, "desired-count", 1, "desired count")

	flags.Int64Var(&ecsCreateCmdTimeout, "timeout", 300, "timeout")

	flags.BoolVarP(&ecsCreateCmdWaitForServiceStable, "service-stable", "w", false, "waits for service to become stable")

}

// RegisterNewTaskDefinition registers task definition
func RegisterNewTaskDefinition(cpu, memory, port int64, service, image string) string {
	taskdefinition := &ecs.TaskDefinition{
		ContainerDefinitions: []ecs.ContainerDefinition{
			ecs.ContainerDefinition{
				Name:              aws.String(service),
				Cpu:               aws.Int64(cpu),
				Memory:            aws.Int64(memory),
				MemoryReservation: aws.Int64(memory),
				PortMappings: []ecs.PortMapping{
					ecs.PortMapping{
						ContainerPort: aws.Int64(port),
						HostPort:      aws.Int64(0),
						Protocol:      "tcp",
					},
				},
				Image: aws.String(image),
			},
		},
		Family: aws.String(service),
	}

	result, err := ecsw.RegisterTaskDefinition(taskdefinition)
	ExitOnError(err, "registering task definition")

	return *result.TaskDefinition.TaskDefinitionArn
}

// CPUAndMemory struct consisting of cpu and memory
type CPUAndMemory struct {
	CPU    int64
	Memory int64
}

// GetCPUAndMemory returns GetCPUAndMemory based on given size
func GetCPUAndMemory(size string) CPUAndMemory {
	sizes := map[string]CPUAndMemory{
		"xsmall":  CPUAndMemory{64, 128},
		"small":   CPUAndMemory{128, 256},
		"medium":  CPUAndMemory{256, 512},
		"large":   CPUAndMemory{512, 1024},
		"xlarge":  CPUAndMemory{1024, 2048},
		"2xlarge": CPUAndMemory{1028, 4096},
	}

	return sizes[size]
}
