package in

import "context"

type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type PublicListModelsResult struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

type PublicListModelsUseCase interface {
	ListModels(ctx context.Context, command PublicListModelsCommand) (*PublicListModelsResult, error)
}
