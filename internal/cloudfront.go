package internal

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

var (
	distributionID string
)

func init() {
	distributionID = os.Getenv("DISTRIBUTION_ID")
}

type CloudFrontInvalidator struct {
	cfClient *cloudfront.Client
}

func (cfi *CloudFrontInvalidator) Invalidate(id string, paths []string) error {
	_, err := cfi.cfClient.CreateInvalidation(context.TODO(), &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(id),
			Paths: &types.Paths{
				Quantity: aws.Int32(int32(len(paths))),
				Items:    paths,
			},
		},
	})

	return err
}

func NewCloudFrontInvalidator() *CloudFrontInvalidator {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	if bucketName == "" {
		log.Fatalf("unable to load bucket name")
	}

	cfClient := cloudfront.NewFromConfig(cfg)

	return &CloudFrontInvalidator{cfClient: cfClient}
}
