package in

type DownloadTrainingDatasetResult struct {
	FieldNames []string
	Data       [][]string
	Filename   string
}

type DownloadTrainingDatasetUseCase interface {
	DownloadTrainingDataset(command DownloadTrainingDatasetCommand) (*DownloadTrainingDatasetResult, error)
}
