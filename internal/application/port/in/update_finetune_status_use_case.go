package in

import "context"

type UpdateFinetuneStatusUseCase interface {
	Execute(ctx context.Context, command UpdateFinetuneStatusCommand) error
}