package web

type AnalyzePromptRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}
