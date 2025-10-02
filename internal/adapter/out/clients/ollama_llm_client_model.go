package clients

type OllamaLLMResponseModel struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	DelayTime     int    `json:"delayTime"`
	ExecutionTime int    `json:"executionTime"`
	Output        []struct {
		Choices []struct {
			Text         string `json:"text"`
			Index        int    `json:"index"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Created         int    `json:"created"`
		ID              string `json:"id"`
		Model           string `json:"model"`
		Object          string `json:"object"`
		SystemFingerprint string `json:"system_fingerprint"`
		Usage           struct {
			CompletionTokens int `json:"completion_tokens"`
			PromptTokens     int `json:"prompt_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	} `json:"output"`
}
