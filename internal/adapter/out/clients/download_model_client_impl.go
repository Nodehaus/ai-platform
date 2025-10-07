package clients

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type DownloadModelClientImpl struct {
	s3Client *s3.Client
	bucket   string
}

func NewDownloadModelClientImpl() (*DownloadModelClientImpl, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if endpointURL := os.Getenv("AWS_ENDPOINT_URL"); endpointURL != "" {
					return aws.Endpoint{
						URL:               endpointURL,
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

	return &DownloadModelClientImpl{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
	}, nil
}

func (c *DownloadModelClientImpl) DownloadModel(ctx context.Context, finetuneID uuid.UUID, modelName string) (io.ReadCloser, int64, error) {
	appEnv := os.Getenv("APP_ENV")
	key := fmt.Sprintf("%s/finetunes/%s/%s.gguf", appEnv, finetuneID.String(), modelName)

	// First get object metadata to check if file exists and get content length
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	headResult, err := c.s3Client.HeadObject(ctx, headInput)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get object metadata for %s: %w", key, err)
	}

	// Get the object for streaming
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	result, err := c.s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get object %s: %w", key, err)
	}

	contentLength := *headResult.ContentLength

	return result.Body, contentLength, nil
}