package clients

type OllamaLLMResponseModel struct {
	Output struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	} `json:"output"`
}
