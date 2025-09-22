package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	portClients "ai-platform/internal/application/port/out/clients"
)

type TrainingDatasetJobClientImpl struct {
	s3Client *s3.Client
	bucket   string
}

func NewTrainingDatasetJobClientImpl() (*TrainingDatasetJobClientImpl, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if endpointURL := os.Getenv("AWS_ENDPOINT_URL"); endpointURL != "" {
					return aws.Endpoint{
						URL:           endpointURL,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	bucket := os.Getenv("APP_S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("APP_S3_BUCKET environment variable is required")
	}

	return &TrainingDatasetJobClientImpl{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
	}, nil
}

func (c *TrainingDatasetJobClientImpl) SubmitJob(ctx context.Context, job portClients.TrainingDatasetJobModel) error {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job to JSON: %w", err)
	}

	key := fmt.Sprintf("jobs/%s_%s.json", time.Now().Format("060102150405"), job.TrainingDatasetID)

	_, err = c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jobJSON),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload job to S3: %w", err)
	}

	return nil
}