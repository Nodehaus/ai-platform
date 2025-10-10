package in

type DownloadDeploymentLogsResult struct {
	FieldNames []string
	Data       [][]string
	Filename   string
}

type DownloadDeploymentLogsUseCase interface {
	DownloadDeploymentLogs(command DownloadDeploymentLogsCommand) (*DownloadDeploymentLogsResult, error)
}
