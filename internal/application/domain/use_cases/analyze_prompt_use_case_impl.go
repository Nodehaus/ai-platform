package use_cases

import (
	"context"
	"fmt"

	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/port/in"
)

type AnalyzePromptUseCaseImpl struct {
	PromptAnalysisService *services.PromptAnalysisService
}

func (uc *AnalyzePromptUseCaseImpl) AnalyzePrompt(command in.AnalyzePromptCommand) (*in.AnalyzePromptResult, error) {
	ctx := context.Background()

	// First call: Get analysis/improvement suggestions
	analysisResult, err := uc.PromptAnalysisService.GetPromptAnalysis(ctx, command.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt analysis: %w", err)
	}

	// Second call: Get JSON structure
	jsonStructure, err := uc.PromptAnalysisService.GetJSONStructure(ctx, command.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON structure: %w", err)
	}

	return &in.AnalyzePromptResult{
		AnalysisResult:          analysisResult,
		JSONObjectFields:        jsonStructure.JSONObjectFields,
		InputField:              jsonStructure.InputField,
		OutputField:             jsonStructure.OutputField,
		ExpectedOutputSizeChars: jsonStructure.ExpectedOutputSizeChars,
	}, nil
}
