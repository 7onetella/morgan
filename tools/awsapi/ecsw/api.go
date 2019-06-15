package ecsw

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

const awsTimeoutDefault = 3

func newECS() (*ecs.ECS, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	cfg.Region = endpoints.UsEast1RegionID

	return ecs.New(cfg), nil
}

func newContextWithTimeout(timeout int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

// ListClusters lists ecs clusters
func ListClusters() (*ecs.ListClustersOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.ListClustersRequest(&ecs.ListClustersInput{})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// DescribeCluster describes ecs cluster
func DescribeClusters(cluster string) (*ecs.DescribeClustersOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.DescribeClustersRequest(&ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// GetClustersForService gets clusters for service
func GetClustersForService(service string) (map[string]string, error) {
	clusters := map[string]string{}

	result, err := ListClusters()
	if err != nil {
		return clusters, err
	}

	for _, c := range result.ClusterArns {
		result2, err := DescribeClusters(c)
		if err != nil {
			return clusters, err
		}

		for _, cdata := range result2.Clusters {
			result3, err := ListServices(*cdata.ClusterName)
			if err != nil {
				return clusters, err
			}

			for _, s := range result3.ServiceArns {

				i := strings.LastIndex(s, "/")
				serviceName := s[i+1:]

				if service == serviceName {
					cluster := *cdata.ClusterName
					clusters[cluster] = serviceName
				}
			}
		}
	}

	return clusters, nil
}

// DescribeServices describes ecs services
func DescribeServices(cluster, service string) (*ecs.DescribeServicesOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.DescribeServicesRequest(&ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []string{service},
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// ListServices lists ecs services
func ListServices(cluster string) (*ecs.ListServicesOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.ListServicesRequest(&ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// UpdateService updates ecs service
func UpdateService(cluster, service, taskdef string, desiredCount int64) (*ecs.UpdateServiceOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.UpdateServiceRequest(&ecs.UpdateServiceInput{
		Cluster:        aws.String(cluster),
		Service:        aws.String(service),
		TaskDefinition: aws.String(taskdef),
		DesiredCount:   aws.Int64(desiredCount),
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// CreateService creates ecs service
func CreateService(cluster, service, taskdef string, desiredCount int64) (*ecs.CreateServiceOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.CreateServiceRequest(&ecs.CreateServiceInput{
		Cluster:        aws.String(cluster),
		ServiceName:    aws.String(service),
		TaskDefinition: aws.String(taskdef),
		DesiredCount:   aws.Int64(desiredCount),
		LaunchType:     ecs.LaunchTypeEc2,
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// TaskDefinitionJSON task definition json
type TaskDefinitionJSON struct {
	Family               string `json:"family"`
	TaskRoleArn          string `json:"taskRoleArn"`
	ExecutionRoleArn     string `json:"executionRoleArn"`
	NetworkMode          string `json:"networkMode"`
	ContainerDefinitions []struct {
		Name                  string `json:"name"`
		Image                 string `json:"image"`
		RepositoryCredentials struct {
			CredentialsParameter string `json:"credentialsParameter"`
		} `json:"repositoryCredentials"`
		CPU               int      `json:"cpu"`
		Memory            int      `json:"memory"`
		MemoryReservation int      `json:"memoryReservation"`
		Links             []string `json:"links"`
		PortMappings      []struct {
			ContainerPort int    `json:"containerPort"`
			HostPort      int    `json:"hostPort"`
			Protocol      string `json:"protocol"`
		} `json:"portMappings"`
		Essential   bool     `json:"essential"`
		EntryPoint  []string `json:"entryPoint"`
		Command     []string `json:"command"`
		Environment []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"environment"`
		MountPoints []struct {
			SourceVolume  string `json:"sourceVolume"`
			ContainerPath string `json:"containerPath"`
			ReadOnly      bool   `json:"readOnly"`
		} `json:"mountPoints"`
		VolumesFrom []struct {
			SourceContainer string `json:"sourceContainer"`
			ReadOnly        bool   `json:"readOnly"`
		} `json:"volumesFrom"`
		LinuxParameters struct {
			Capabilities struct {
				Add  []string `json:"add"`
				Drop []string `json:"drop"`
			} `json:"capabilities"`
			Devices []struct {
				HostPath      string   `json:"hostPath"`
				ContainerPath string   `json:"containerPath"`
				Permissions   []string `json:"permissions"`
			} `json:"devices"`
			InitProcessEnabled bool `json:"initProcessEnabled"`
			SharedMemorySize   int  `json:"sharedMemorySize"`
			Tmpfs              []struct {
				ContainerPath string   `json:"containerPath"`
				Size          int      `json:"size"`
				MountOptions  []string `json:"mountOptions"`
			} `json:"tmpfs"`
		} `json:"linuxParameters"`
		Secrets []struct {
			Name      string `json:"name"`
			ValueFrom string `json:"valueFrom"`
		} `json:"secrets"`
		DependsOn []struct {
			ContainerName string `json:"containerName"`
			Condition     string `json:"condition"`
		} `json:"dependsOn"`
		StartTimeout           int      `json:"startTimeout"`
		StopTimeout            int      `json:"stopTimeout"`
		Hostname               string   `json:"hostname"`
		User                   string   `json:"user"`
		WorkingDirectory       string   `json:"workingDirectory"`
		DisableNetworking      bool     `json:"disableNetworking"`
		Privileged             bool     `json:"privileged"`
		ReadonlyRootFilesystem bool     `json:"readonlyRootFilesystem"`
		DNSServers             []string `json:"dnsServers"`
		DNSSearchDomains       []string `json:"dnsSearchDomains"`
		ExtraHosts             []struct {
			Hostname  string `json:"hostname"`
			IPAddress string `json:"ipAddress"`
		} `json:"extraHosts"`
		DockerSecurityOptions []string `json:"dockerSecurityOptions"`
		Interactive           bool     `json:"interactive"`
		PseudoTerminal        bool     `json:"pseudoTerminal"`
		DockerLabels          struct {
			KeyName string `json:"KeyName"`
		} `json:"dockerLabels"`
		Ulimits []struct {
			Name      string `json:"name"`
			SoftLimit int    `json:"softLimit"`
			HardLimit int    `json:"hardLimit"`
		} `json:"ulimits"`
		LogConfiguration struct {
			LogDriver string `json:"logDriver"`
			Options   struct {
				KeyName string `json:"KeyName"`
			} `json:"options"`
			SecretOptions []struct {
				Name      string `json:"name"`
				ValueFrom string `json:"valueFrom"`
			} `json:"secretOptions"`
		} `json:"logConfiguration"`
		HealthCheck struct {
			Command     []string `json:"command"`
			Interval    int      `json:"interval"`
			Timeout     int      `json:"timeout"`
			Retries     int      `json:"retries"`
			StartPeriod int      `json:"startPeriod"`
		} `json:"healthCheck"`
		SystemControls []struct {
			Namespace string `json:"namespace"`
			Value     string `json:"value"`
		} `json:"systemControls"`
		ResourceRequirements []struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"resourceRequirements"`
	} `json:"containerDefinitions"`
	Volumes []struct {
		Name string `json:"name"`
		Host struct {
			SourcePath string `json:"sourcePath"`
		} `json:"host"`
		DockerVolumeConfiguration struct {
			Scope         string `json:"scope"`
			Autoprovision bool   `json:"autoprovision"`
			Driver        string `json:"driver"`
			DriverOpts    struct {
				KeyName string `json:"KeyName"`
			} `json:"driverOpts"`
			Labels struct {
				KeyName string `json:"KeyName"`
			} `json:"labels"`
		} `json:"dockerVolumeConfiguration"`
	} `json:"volumes"`
	PlacementConstraints []struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	} `json:"placementConstraints"`
	RequiresCompatibilities []string `json:"requiresCompatibilities"`
	CPU                     string   `json:"cpu"`
	Memory                  string   `json:"memory"`
	Tags                    []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
	PidMode            string `json:"pidMode"`
	IpcMode            string `json:"ipcMode"`
	ProxyConfiguration struct {
		Type          string `json:"type"`
		ContainerName string `json:"containerName"`
		Properties    []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"properties"`
	} `json:"proxyConfiguration"`
}

// ServiceStable waits for ecs service to be stable
func ServiceStable(cluster, service string, timeout int64) error {
	svc, err := newECS()
	if err != nil {
		return err
	}

	ctx, cancel := newContextWithTimeout(timeout)
	defer cancel()

	err = svc.WaitUntilServicesStable(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []string{service},
	})

	return err
}

// DescribeTaskDefinition desribes task definition
func DescribeTaskDefinition(taskdef string) (*ecs.DescribeTaskDefinitionOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.DescribeTaskDefinitionRequest(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskdef),
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// RegisterTaskDefinition registers task definition
func RegisterTaskDefinition(td *ecs.TaskDefinition) (*ecs.RegisterTaskDefinitionOutput, error) {
	svc, err := newECS()
	if err != nil {
		return nil, err
	}

	req := svc.RegisterTaskDefinitionRequest(&ecs.RegisterTaskDefinitionInput{
		Family:                  td.Family,
		ContainerDefinitions:    td.ContainerDefinitions,
		IpcMode:                 td.IpcMode,
		Memory:                  td.Memory,
		NetworkMode:             td.NetworkMode,
		PidMode:                 td.PidMode,
		PlacementConstraints:    td.PlacementConstraints,
		RequiresCompatibilities: td.RequiresCompatibilities,
		Volumes:                 td.Volumes,
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}
