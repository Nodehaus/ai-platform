package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"ai-platform/internal/application/domain/entities"
	portClients "ai-platform/internal/application/port/out/clients"
)

type TrainingDatasetResultsClientImpl struct {
	s3Client *s3.Client
	bucket   string
}

func NewTrainingDatasetResultsClientImpl() (*TrainingDatasetResultsClientImpl, error) {
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

	return &TrainingDatasetResultsClientImpl{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
	}, nil
}

func (c *TrainingDatasetResultsClientImpl) GetTrainingDatasetResults(ctx context.Context, trainingDatasetID uuid.UUID, fieldNames []string) (*portClients.TrainingDatasetResult, error) {
	// List all JSON files in the datasets/{training_dataset_id}/ path
	prefix := fmt.Sprintf("datasets/%s/", trainingDatasetID.String())

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	}

	var allObjects []types.Object
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, listInput)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}
		allObjects = append(allObjects, page.Contents...)
	}

	// Filter for JSON files only
	var jsonFiles []types.Object
	for _, obj := range allObjects {
		if len(*obj.Key) > 5 && (*obj.Key)[len(*obj.Key)-5:] == ".json" {
			jsonFiles = append(jsonFiles, obj)
		}
	}

	if len(jsonFiles) == 0 {
		return nil, fmt.Errorf("no JSON files found for training dataset %s", trainingDatasetID.String())
	}

	var totalGenerationTime float64
	var allAnnotations []AnnotationModel

	// Process each JSON file
	for _, file := range jsonFiles {
		fileData, err := c.downloadFile(ctx, *file.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to download file %s: %w", *file.Key, err)
		}

		// Parse as raw JSON first to handle dynamic fields
		var rawData map[string]interface{}
		if err := json.Unmarshal(fileData, &rawData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file %s: %w", *file.Key, err)
		}

		// Extract the structured fields
		fileModel := TrainingDatasetResultsFileModel{}
		if totalGenTime, ok := rawData["total_generation_time"].(float64); ok {
			fileModel.TotalGenerationTime = totalGenTime
		}

		// Parse annotations with dynamic fields
		if annotationsRaw, ok := rawData["annotations"].([]interface{}); ok {
			for _, annotationRaw := range annotationsRaw {
				if annotationMap, ok := annotationRaw.(map[string]interface{}); ok {
					annotation := AnnotationModel{
						Fields: make(map[string]interface{}),
					}

					// Extract fixed fields
					if docID, ok := annotationMap["document_id"].(string); ok {
						annotation.DocumentID = docID
					}
					if withinStart, ok := annotationMap["within_start"].(float64); ok {
						annotation.WithinStart = int(withinStart)
					}
					if withinEnd, ok := annotationMap["within_end"].(float64); ok {
						annotation.WithinEnd = int(withinEnd)
					}
					if inferenceTime, ok := annotationMap["inference_time_seconds"].(float64); ok {
						annotation.InferenceTimeSeconds = inferenceTime
					}

					// Store all fields for dynamic access
					for key, value := range annotationMap {
						annotation.Fields[key] = value
					}

					fileModel.Annotations = append(fileModel.Annotations, annotation)
				}
			}
		}

		totalGenerationTime += fileModel.TotalGenerationTime
		allAnnotations = append(allAnnotations, fileModel.Annotations...)
	}

	// Convert annotations to TrainingDataItem entities
	trainingDataItems, err := c.convertAnnotationsToTrainingDataItems(allAnnotations, fieldNames)
	if err != nil {
		return nil, fmt.Errorf("failed to convert annotations: %w", err)
	}

	return &portClients.TrainingDatasetResult{
		TotalGenerationTimeSeconds: totalGenerationTime,
		TrainingDataItems:          trainingDataItems,
	}, nil
}

func (c *TrainingDatasetResultsClientImpl) downloadFile(ctx context.Context, key string) ([]byte, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	result, err := c.s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (c *TrainingDatasetResultsClientImpl) convertAnnotationsToTrainingDataItems(annotations []AnnotationModel, fieldNames []string) ([]entities.TrainingDataItem, error) {
	var trainingDataItems []entities.TrainingDataItem

	for _, annotation := range annotations {
		// Extract values in the order specified by fieldNames
		values := make([]string, len(fieldNames))

		// Map field names to values in the correct order
		for i, fieldName := range fieldNames {
			if value, exists := annotation.Fields[fieldName]; exists {
				values[i] = fmt.Sprintf("%v", value)
			} else {
				return nil, fmt.Errorf("field %s not found in annotation", fieldName)
			}
		}

		// Create TrainingDataItem
		item := entities.TrainingDataItem{
			ID:                    uuid.New(),
			Values:                values,
			SourceDocument:        &annotation.DocumentID,
			SourceDocumentStart:   func() *string { s := fmt.Sprintf("%d", annotation.WithinStart); return &s }(),
			SourceDocumentEnd:     func() *string { s := fmt.Sprintf("%d", annotation.WithinEnd); return &s }(),
			GenerationTimeSeconds: annotation.InferenceTimeSeconds,
			Deleted:               false,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}

		trainingDataItems = append(trainingDataItems, item)
	}

	return trainingDataItems, nil
}