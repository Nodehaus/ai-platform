package web

import (
	"ai-platform/internal/application/domain/entities"
	"time"

	"github.com/google/uuid"
)

type DeploymentLogSample struct {
	CreatedAt time.Time `json:"created_at"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
}

type GetDeploymentResponse struct {
	ID         uuid.UUID             `json:"id"`
	ModelName  string                `json:"model_name"`
	APIKey     string                `json:"api_key"`
	ProjectID  uuid.UUID             `json:"project_id"`
	FinetuneID *uuid.UUID            `json:"finetune_id"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
	LogsSample []DeploymentLogSample `json:"logs_sample"`
}

func NewGetDeploymentResponse(deployment *entities.Deployment, logs []*entities.DeploymentLogs) *GetDeploymentResponse {
	logsSample := make([]DeploymentLogSample, 0, len(logs))
	for _, log := range logs {
		logsSample = append(logsSample, DeploymentLogSample{
			CreatedAt: log.CreatedAt,
			Input:     log.Input,
			Output:    log.Output,
		})
	}

	return &GetDeploymentResponse{
		ID:         deployment.ID,
		ModelName:  deployment.ModelName,
		APIKey:     deployment.APIKey,
		ProjectID:  deployment.ProjectID,
		FinetuneID: deployment.FinetuneID,
		CreatedAt:  deployment.CreatedAt,
		UpdatedAt:  deployment.UpdatedAt,
		LogsSample: logsSample,
	}
}
