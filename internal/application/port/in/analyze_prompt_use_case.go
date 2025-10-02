package in

type AnalyzePromptResult struct {
	AnalysisResult          string
	JSONObjectFields        map[string]string
	InputField              string
	OutputField             string
	ExpectedOutputSizeChars int
}

type AnalyzePromptUseCase interface {
	AnalyzePrompt(command AnalyzePromptCommand) (*AnalyzePromptResult, error)
}
