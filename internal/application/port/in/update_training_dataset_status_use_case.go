package in

import "context"

type UpdateTrainingDatasetStatusUseCase interface {
	Execute(ctx context.Context, command UpdateTrainingDatasetStatusCommand) error
}