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
	"os"
	"strconv"
	"strings"

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ecsDescribeCmdCluster string
var ecsDescribeCmdTimeout int64

var ecsDescribeCmd = &cobra.Command{
	Use:     "describe-services <service names>",
	Short:   "Describes ecs",
	Long:    `Describes ecs`,
	Example: "foo-svc 1.0.0 --cluster api-cluster",
	Aliases: []string{"describe"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cluster := ecsDescribeCmdCluster
		service := args[0]

		// if cluster is not specified, then assume there is only one cluster and use that cluster
		if len(cluster) == 0 {
			clusters := GetClustersForService(service)

			CheckForClusterAmbiguity(clusters)

			cluster = GetClusterForService(clusters, service)
		}

		services := args[0:]

		result, err := ecsw.DescribeServices(cluster, services...)
		ExitOnError(err, "describing services")
		if len(result.Services) == 0 {
			ExitOnError(errors.New("search result count 0"), "finding services")
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Pending", "Running", "Desired", "TaskDef"})

		for _, s := range result.Services {
			taskdefarn := *s.TaskDefinition
			lastSlash := strings.LastIndex(taskdefarn, "/")
			taskdef := taskdefarn[lastSlash+1:]
			toString := func(num *int64) string {
				v := strconv.Itoa(int(*num))
				return v
			}
			table.Append([]string{*s.ServiceName, toString(s.PendingCount), toString(s.RunningCount), toString(s.DesiredCount), taskdef})
		}
		table.Render()

	},
}

func init() {

	ecsCmd.AddCommand(ecsDescribeCmd)

	flags := ecsDescribeCmd.Flags()

	flags.StringVar(&ecsDescribeCmdCluster, "cluster", "", "ecs cluster")

}
