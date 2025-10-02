package web

type AnalyzePromptResponse struct {
	AnalysisResult          string            `json:"analysis_result"`
	JSONObjectFields        map[string]string `json:"json_object_fields"`
	InputField              string            `json:"input_field"`
	OutputField             string            `json:"output_field"`
	ExpectedOutputSizeChars int               `json:"expected_output_size_chars"`
}

func NewAnalyzePromptResponse(
	analysisResult string,
	jsonObjectFields map[string]string,
	inputField string,
	outputField string,
	expectedOutputSizeChars int,
) *AnalyzePromptResponse {
	return &AnalyzePromptResponse{
		AnalysisResult:          analysisResult,
		JSONObjectFields:        jsonObjectFields,
		InputField:              inputField,
		OutputField:             outputField,
		ExpectedOutputSizeChars: expectedOutputSizeChars,
	}
}
