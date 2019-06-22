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

	"github.com/7onetella/morgan/tools/awsapi/ecsw"
	"github.com/spf13/cobra"
)

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Automation for ecs",
	Long:  `Automation for ecs`,
}

func init() {
	awsCmd.AddCommand(ecsCmd)
}

// CheckForClusterAmbiguity checks for existence of given service in all clusters
func CheckForClusterAmbiguity(clusters map[string]string) {
	if len(clusters) > 1 {
		Failure("checking for cluster ambiquity")
		Info("more than one clusters with same service name found. you must explicitly specify --cluster argument")
		Newline()
		for k := range clusters {
			Print(indentation + "  - " + k + "\n")
		}
		os.Exit(1)
	}
}

// GetClusterForService gets the cluster for given service
func GetClusterForService(clusters map[string]string, service string) string {

	if len(clusters) == 1 {
		for cluster, svc := range clusters {
			if svc == service {
				return cluster
			}
		}
	}

	return ""
}

// GetClustersForService gets the clusters for given service
func GetClustersForService(service string) map[string]string {
	clusters, err := ecsw.GetClustersForService(service)
	ExitOnError(err, "getting clusters for service")

	return clusters
}
