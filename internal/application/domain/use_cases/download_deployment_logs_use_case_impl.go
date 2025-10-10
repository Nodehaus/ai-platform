package use_cases

import (
	"fmt"

	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
	"ai-platform/internal/application/port/out/persistence"
)

type DownloadDeploymentLogsUseCaseImpl struct {
	DeploymentLogsRepository persistence.DeploymentLogsRepository
	DeploymentRepository     persistence.DeploymentRepository
	DeploymentService        *services.DeploymentService
}

func (uc *DownloadDeploymentLogsUseCaseImpl) DownloadDeploymentLogs(command in.DownloadDeploymentLogsCommand) (*in.DownloadDeploymentLogsResult, error) {
	// Validate project access
	err := uc.DeploymentService.ValidateProjectAccess(command.ProjectID, command.OwnerID)
	if err != nil {
		return nil, err
	}

	// Get deployment to verify it exists and belongs to the project
	deployment, err := uc.DeploymentRepository.GetByID(command.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	if deployment == nil {
		return nil, fmt.Errorf("deployment not found")
	}

	// Verify deployment belongs to project
	if deployment.ProjectID != command.ProjectID {
		return nil, fmt.Errorf("deployment does not belong to this project")
	}

	// Get all deployment logs
	logs, err := uc.DeploymentLogsRepository.GetAll(command.DeploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment logs: %w", err)
	}

	// Convert logs to CSV format
	fieldNames := []string{"date", "input", "output"}
	var data [][]string
	for _, log := range logs {
		row := []string{
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.Input,
			log.Output,
		}
		data = append(data, row)
	}

	// Generate filename: deployment_logs_{model_name}.csv
	filename := fmt.Sprintf("deployment_logs_%s.csv", deployment.ModelName)

	return &in.DownloadDeploymentLogsResult{
		FieldNames: fieldNames,
		Data:       data,
		Filename:   filename,
	}, nil
}
