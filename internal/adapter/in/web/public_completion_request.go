package web

type PublicCompletionRequest struct {
	Model       string  `json:"model" binding:"required"`
	Prompt      string  `json:"prompt" binding:"required"`
	MaxTokens   *int    `json:"max_tokens"`
	Temperature *float64 `json:"temperature"`
	TopP        *float64 `json:"top_p"`
	Stream      *bool    `json:"stream"`
}
