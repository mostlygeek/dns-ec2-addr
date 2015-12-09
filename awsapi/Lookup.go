package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetIPs(region, tagApp, tagType string) ([]string, error) {
	ips := make([]string, 0)
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))

	params := &ec2.DescribeInstancesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{ // Required
				Name: aws.String("tag:App"),
				Values: []*string{
					aws.String(tagApp), // Required
				},
			},
			{ // Required
				Name: aws.String("tag:Type"),
				Values: []*string{
					aws.String(tagType), // Required
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return nil, err
	}

	// Pretty-print the response data.
	for _, reservations := range resp.Reservations {
		for _, inst := range reservations.Instances {
			for _, iface := range inst.NetworkInterfaces {
				if iface.Association.PublicIp != nil {
					ips = append(ips, *iface.Association.PublicIp)
				}
			}
		}
	}

	return ips, nil
}
