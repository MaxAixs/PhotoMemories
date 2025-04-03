package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
	region     string
}

func NewS3Client(region string, bucketName string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &Client{
		s3Client:   client,
		bucketName: bucketName,
		region:     region,
	}, nil
}
