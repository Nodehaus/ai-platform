package in

type UploadTrainingDatasetResult struct {
	ItemsAdded int
	TotalItems int
}

type UploadTrainingDatasetUseCase interface {
	UploadTrainingDataset(command UploadTrainingDatasetCommand) (*UploadTrainingDatasetResult, error)
}
