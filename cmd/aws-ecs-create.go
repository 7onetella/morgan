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
var ecsCreateCmdEnvVars []string
var ecsCreateCmdDesiredCount int64
var ecsCreateCmdTimeout int64
var ecsCreateCmdWaitForServiceStable bool

var ecsCreateCmd = &cobra.Command{
	Use:   "create <service-name> <size> <port> <docker-image>",
	Short: "Creates ecs",
	Long: `
Presently, there is no support for attaching the ecs service to ALB. ALB will be supported shortly. 
You are expected to have some other mechanism to put your service behind reverse proxy. 
Consul and Fabio proxy can be used to attach standalone service such as our example.

The value of <size> parameter is t-shirt sized.
* xsmall  : CPU 63,   Memory 128
* small   : CPU 128,  Memory 256
* medium  : CPU 256,  Memory 512
* large   : CPU 512,  Memory 1024
* xlarge  : CPU 1024, Memory 2048
* 2xlarge : CPU 2024, Memory 4096`,
	Example: `hello-world xsmall 8080 7onetealla/ref-api:latest \
	--cluster Development \
	-e NAME=web \
	-e URLPREFIX=foo-svc.example.com/`,
	Aliases: []string{"create-service"},
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

		envs := ConvertKeyValuePairArgSliceToMap(ecsCreateCmdEnvVars)

		if len(taskdef) == 0 {
			cm := GetCPUAndMemory(size)
			taskdef = RegisterNewTaskDefinition(cm.CPU, cm.Memory, int64(port), service, image, envs)
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

	flags.StringVarP(&ecsCreateCmdCluster, "cluster", "c", "", "optional: ecs cluster name")

	flags.StringVarP(&ecsCreateCmdService, "service", "s", "", "required: service name")

	flags.StringVar(&ecsCreateCmdTaskDefinition, "task-definition", "", "optional: task definition. e.g. family:revision")

	flags.Int64Var(&ecsCreateCmdDesiredCount, "desired-count", 1, "optional: desired count")

	flags.BoolVarP(&ecsCreateCmdWaitForServiceStable, "service-stable", "w", false, "optional: waits for service to become stable")

	flags.Int64Var(&ecsCreateCmdTimeout, "timeout", 300, "optional: timeout for service stable")

	flags.StringSliceVarP(&ecsCreateCmdEnvVars, "env", "e", []string{}, "optional: environment variables. e.g. -e key=value")

}

// RegisterNewTaskDefinition registers task definition
func RegisterNewTaskDefinition(cpu, memory, port int64, service, image string, environmentVars map[string]string) string {
	envs := []ecs.KeyValuePair{}

	for k, v := range environmentVars {
		envs = append(envs, ecs.KeyValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}

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
				Image:       aws.String(image),
				Environment: envs,
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
