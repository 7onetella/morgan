package ecsw

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
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
