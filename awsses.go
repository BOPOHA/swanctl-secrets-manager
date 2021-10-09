package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"log"
)

func getSesClient() *sesv2.Client {
	// https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/ec2-imds/
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cliImds := imds.NewFromConfig(cfg)
	region, err := cliImds.GetRegion(context.TODO(), &imds.GetRegionInput{})
	if err != nil {
		log.Fatalf("Unable to retrieve the region from the EC2 instance %v\n", err)
	}
	cfg.Region = region.Region

	return sesv2.NewFromConfig(cfg)

}
