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

// DescribeInstanceByTagAndValue describes instance by tag and tag values
func DescribeInstanceByTagAndValue(tagName string, values ...string) (*ec2.DescribeInstancesOutput, error) {
	svc, err := newEC2()
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeInstancesInput{}
	if len(tagName) > 0 && len(values) > 0 {
		input.Filters = []ec2.Filter{
			ec2.Filter{
				Name:   aws.String("tag:" + tagName),
				Values: values,
			},
		}
	}

	req := svc.DescribeInstancesRequest(input)

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// DescribeInstanceByNameTag describes instance by name tag
func DescribeInstanceByNameTag(name string) (*ec2.DescribeInstancesOutput, error) {
	return DescribeInstanceByTagAndValue("Name", name)
}

// GetInstanceIDsByNames gets instance ids by names
func GetInstanceIDsByNames(names []string) (map[string]string, error) {

	instanceIDs := map[string]string{}

	for _, name := range names {
		resp, err := DescribeInstanceByNameTag(name)
		if err != nil {
			return map[string]string{}, err
		}
		if len(resp.Reservations) == 0 {
			return map[string]string{}, fmt.Errorf("can not find instance for %s", name)
		}
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				instanceIDs[*instance.InstanceId] = name
			}
		}
	}

	return instanceIDs, nil
}
