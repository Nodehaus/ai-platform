package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrainingDatasetService_GenerateCsvFilename(t *testing.T) {
	service := &TrainingDatasetService{}

	tests := []struct {
		name        string
		projectName string
		version     int
		expected    string
	}{
		{
			name:        "Simple project name",
			projectName: "My Project",
			version:     1,
			expected:    "dataset_my_project_v1.csv",
		},
		{
			name:        "Project name with special characters",
			projectName: "My-Project@2024!",
			version:     2,
			expected:    "dataset_myproject2024_v2.csv",
		},
		{
			name:        "Project name with multiple spaces",
			projectName: "My  Test  Project",
			version:     3,
			expected:    "dataset_my_test_project_v3.csv",
		},
		{
			name:        "Project name with underscores",
			projectName: "My_Project_Name",
			version:     1,
			expected:    "dataset_my_project_name_v1.csv",
		},
		{
			name:        "Project name with mixed case and special chars",
			projectName: "AI-Platform (2024) & ML Tools!",
			version:     5,
			expected:    "dataset_aiplatform_2024_ml_tools_v5.csv",
		},
		{
			name:        "Project name with leading/trailing spaces",
			projectName: "  My Project  ",
			version:     1,
			expected:    "dataset_my_project_v1.csv",
		},
		{
			name:        "Project name all uppercase",
			projectName: "MYPROJECT",
			version:     10,
			expected:    "dataset_myproject_v10.csv",
		},
		{
			name:        "Project name with numbers",
			projectName: "Project 123",
			version:     1,
			expected:    "dataset_project_123_v1.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateCsvFilename(tt.projectName, tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
