package ec2w

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const awsTimeoutDefault = 3

func newEC2() (*ec2.EC2, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	cfg.Region = endpoints.UsEast1RegionID

	return ec2.New(cfg), nil
}

func newContextWithTimeout(timeout int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

// StartInstances starts instances
func StartInstances(instanceIDs []string) (*ec2.StartInstancesOutput, error) {
	svc, err := newEC2()
	if err != nil {
		return nil, err
	}

	req := svc.StartInstancesRequest(&ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// StopInstances stops instances
func StopInstances(instanceIDs []string) (*ec2.StopInstancesOutput, error) {
	svc, err := newEC2()
	if err != nil {
		return nil, err
	}

	req := svc.StopInstancesRequest(&ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// DescribeInstanceByNameTag describes instance by name tag
func DescribeInstanceByNameTag(name string) (*ec2.DescribeInstancesOutput, error) {
	svc, err := newEC2()
	if err != nil {
		return nil, err
	}

	req := svc.DescribeInstancesRequest(&ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			ec2.Filter{
				Name:   aws.String("tag:Name"),
				Values: []string{name},
			},
		},
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// GetInstanceIDsByNames gets instance ids by names
func GetInstanceIDsByNames(names []string) ([]string, error) {

	instanceIDs := []string{}

	for _, name := range names {
		resp, err := DescribeInstanceByNameTag(name)
		if err != nil {
			return []string{}, err
		}
		if len(resp.Reservations) == 0 {
			return []string{}, fmt.Errorf("can not find instance for %s", name)
		}
		instanceIDs = append(instanceIDs, *resp.Reservations[0].Instances[0].InstanceId)
	}

	return instanceIDs, nil
}
