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
	Use:     "describe <service names>",
	Short:   "Describes ecs",
	Long:    `Describes ecs`,
	Example: "foo-svc 1.0.0 --cluster api-cluster",
	Aliases: []string{"describe-services"},
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		specifiedCluster := ecsDescribeCmdCluster
		clusterMembers := map[string][]string{}

		isClusterSpecified := len(specifiedCluster) > 0
		isServicesSpecified := len(args) > 0

		// if no services are specified then get all services from all clusters
		if !isServicesSpecified {
			var clusters []string
			if !isClusterSpecified {
				result, err := ecsw.ListClusters()
				ExitOnError(err, "listing clusters")
				for _, clusterARN := range result.ClusterArns {
					slashIndex := strings.LastIndex(clusterARN, "/")
					cluster := clusterARN[slashIndex+1:]
					clusters = append(clusters, cluster)
				}
			} else {
				// if cluster is specified, add specified cluster to clusters
				clusters = []string{specifiedCluster}
			}

			for _, cluster := range clusters {
				result, err := ecsw.ListServices(cluster)
				ExitOnError(err, "listing services")
				serviceARNs := result.ServiceArns
				if len(serviceARNs) == 0 {
					continue
				}
				clusterMembers[cluster] = serviceARNs
			}
		}

		if isServicesSpecified {
			if !isClusterSpecified {
				for _, service := range args[0:] {
					clustersForSvc, err := ecsw.GetClustersForService(service)
					ExitOnError(err, "getting clusters for service")
					for cluster := range clustersForSvc {
						// pull services for cluster
						currServices := clusterMembers[cluster]
						// if current services does not contain service
						var isServiceFound bool
						for _, serivceName := range currServices {
							if serivceName == service {
								isServiceFound = true
								break
							}
						}
						if !isServiceFound {
							currServices = append(currServices, service)
						}
						// put modified services back into clusterMembers
						clusterMembers[cluster] = currServices
					}
				}
			} else {
				clusterMembers[specifiedCluster] = args[0:]
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Cluster", "Name", "Pending", "Running", "Desired", "TaskDef", "Tags"})

		for cluster, serviceARNS := range clusterMembers {
			result, err := ecsw.DescribeServices(cluster, serviceARNS...)
			ExitOnError(err, "describing services")
			if len(result.Services) == 0 {
				ExitOnError(errors.New("search result count 0"), "finding services")
			}

			for _, s := range result.Services {
				taskdef := parseTaskDefinitionStr(*s.TaskDefinition)
				tags := getTags(taskdef)
				table.Append([]string{cluster, *s.ServiceName, toString(s.PendingCount), toString(s.RunningCount), toString(s.DesiredCount), taskdef, strings.Join(tags, ",")})
			}
		}

		table.Render()

	},
}

func init() {

	ecsCmd.AddCommand(ecsDescribeCmd)

	flags := ecsDescribeCmd.Flags()

	flags.StringVar(&ecsDescribeCmdCluster, "cluster", "", "ecs cluster")

}

func parseTaskDefinitionStr(taskdefARN string) string {
	lastSlash := strings.LastIndex(taskdefARN, "/")
	taskdef := taskdefARN[lastSlash+1:]
	return taskdef
}

func toString(num *int64) string {
	v := strconv.Itoa(int(*num))
	return v
}

func getTags(taskdef string) []string {
	var tags []string
	t, err := ecsw.DescribeTaskDefinition(taskdef)
	if err == nil {
		for _, cd := range t.TaskDefinition.ContainerDefinitions {
			image := *cd.Image
			colonIndex := strings.LastIndex(image, ":")
			tag := image[colonIndex+1:]
			tags = append(tags, tag)
		}
	}
	return tags
}
