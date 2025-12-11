package clients

type StreamChunk struct {
	Delta struct {
		Role    string `json:"role,omitempty"`
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type StreamMetadata struct {
	TokensIn      int
	TokensOut     int
	DelayTime     int
	ExecutionTime int
}
