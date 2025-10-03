package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai-platform/internal/application/port/out/clients"
)

type PromptAnalysisService struct {
	OllamaLLMClient clients.OllamaLLMClient
}

type JSONStructureResponse struct {
	JSONObjectFields        map[string]string `json:"json_object_fields"`
	InputField              string            `json:"input_field"`
	OutputField             string            `json:"output_field"`
	ExpectedOutputSizeChars int               `json:"expected_output_size_chars"`
}

// GetPromptAnalysis calls the LLM to analyze and suggest improvements for the prompt
func (s *PromptAnalysisService) GetPromptAnalysis(ctx context.Context, userPrompt string) (string, error) {
	promptTemplate := `Analyze the following use case description for a prompt that will be used to create a training dataset for LLM fine-tuning. Provide suggestions on how to improve it to generate better training data. Focus on clarity, specificity, and structure. Keep the analysis short and focus on the most important changes in not more than 4 sentences. Suggest a rewrite of the prompt but do not add any information about how many items to output or what structure to output, we will add the information about the output size and structure separately.

PROMPT:
%s

ANALYSIS:`

	llmPrompt := fmt.Sprintf(promptTemplate, userPrompt)

	response, err := s.OllamaLLMClient.GenerateCompletion(ctx, llmPrompt, 500, 0.7, 0.9)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return strings.TrimSpace(response), nil
}

// GetJSONStructure calls the LLM to determine the JSON object structure for training data
func (s *PromptAnalysisService) GetJSONStructure(ctx context.Context, userPrompt string) (*JSONStructureResponse, error) {
	promptTemplate := `Tell me what list of JSON objects an LLM would output for a given prompt. It does not need to output any source text or reference, **ignore** any mentions of source text. Do not add examples. Include minimum 2 fields, maximum 4 fields. We want to use the output data for LLM training. Decide which fields we will use as input and output of the training. Also return the expected size of the output in characters based on your output field description. Use the following output format:

` + "```json" + `
{
  "json_object_fields": {
    "field_name": "Description of what this field represents",
    "another_field": "Description of this field",
    "yet_another_field": "More descriptions"
  },
  "input_field": "field_name",
  "output_field": "another_field",
  "expected_output_size_chars": a_reasonable_int_for_output_size
}
` + "```" + `

PROMPT:

%s

JSON:`

	llmPrompt := fmt.Sprintf(promptTemplate, userPrompt)

	response, err := s.OllamaLLMClient.GenerateCompletion(ctx, llmPrompt, 800, 0.5, 0.9)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Parse the JSON response
	jsonResponse, err := s.extractJSON(response)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON from response: %w", err)
	}

	var result JSONStructureResponse
	if err := json.Unmarshal([]byte(jsonResponse), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}

// extractJSON attempts to extract JSON from the LLM response
func (s *PromptAnalysisService) extractJSON(response string) (string, error) {
	// Try to find JSON between ```json and ``` markers
	if strings.Contains(response, "```json") {
		start := strings.Index(response, "```json") + 7
		end := strings.Index(response[start:], "```")
		if end != -1 {
			return strings.TrimSpace(response[start : start+end]), nil
		}
	}

	// Try to find JSON between ``` markers
	if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		end := strings.Index(response[start:], "```")
		if end != -1 {
			return strings.TrimSpace(response[start : start+end]), nil
		}
	}

	// Try to find JSON directly by looking for { and }
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(response[start : end+1]), nil
	}

	return "", fmt.Errorf("could not find valid JSON in response")
}
