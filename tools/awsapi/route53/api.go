package route53

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go/aws"
)

const awsTimeoutDefault = 3

func newRoute53() (*route53.Route53, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	cfg.Region = endpoints.UsEast1RegionID

	return route53.New(cfg), nil
}

func newContextWithTimeout(timeout int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

// ListHostedZones lists hosted zones
func ListHostedZones() (*route53.ListHostedZonesOutput, error) {
	svc, err := newRoute53()
	if err != nil {
		return nil, err
	}

	req := svc.ListHostedZonesRequest(&route53.ListHostedZonesInput{})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}

// ARecordUpsert upsert A record
func ARecordUpsert(name, value, hostedZoneID string) (*route53.ChangeResourceRecordSetsOutput, error) {
	svc, err := newRoute53()
	if err != nil {
		return nil, err
	}

	req := svc.ChangeResourceRecordSetsRequest(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []route53.Change{
				route53.Change{
					Action: route53.ChangeActionUpsert,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(name),
						ResourceRecords: []route53.ResourceRecord{
							route53.ResourceRecord{
								Value: aws.String(value),
							},
						},
						Type: route53.RRTypeA,
						TTL:  aws.Int64(300),
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneID),
	})

	ctx, cancel := newContextWithTimeout(awsTimeoutDefault)
	defer cancel()

	return req.Send(ctx)
}
