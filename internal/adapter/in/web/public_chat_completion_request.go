package web

type PublicChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PublicChatCompletionRequest struct {
	Model       string               `json:"model" binding:"required"`
	Messages    []PublicChatMessage  `json:"messages" binding:"required,min=1"`
	MaxTokens   *int                 `json:"max_tokens"`
	Temperature *float64             `json:"temperature"`
	TopP        *float64             `json:"top_p"`
	Stream      *bool                `json:"stream"`
}
