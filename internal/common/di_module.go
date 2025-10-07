package common

import (
	"os"

	"go.uber.org/fx"

	"ai-platform/internal/adapter/in/web"
	"ai-platform/internal/adapter/out/clients"
	"ai-platform/internal/adapter/out/persistence"
	"ai-platform/internal/application/domain/services"
	"ai-platform/internal/application/domain/use_cases"
	"ai-platform/internal/application/port/in"
	clientsPort "ai-platform/internal/application/port/out/clients"
	persistencePort "ai-platform/internal/application/port/out/persistence"
	"ai-platform/internal/database"
	"ai-platform/internal/server"
)

func NewUserRepository(dbService database.Service) persistencePort.UserRepository {
	return &persistence.UserRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewProjectRepository(dbService database.Service) persistencePort.ProjectRepository {
	return &persistence.ProjectRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewTrainingDatasetRepository(dbService database.Service) persistencePort.TrainingDatasetRepository {
	return &persistence.TrainingDatasetRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewCorpusRepository(dbService database.Service) persistencePort.CorpusRepository {
	return &persistence.CorpusRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewPromptRepository(dbService database.Service) persistencePort.PromptRepository {
	return &persistence.PromptRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewFinetuneRepository(dbService database.Service) persistencePort.FinetuneRepository {
	return &persistence.FinetuneRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewDeploymentRepository(dbService database.Service) persistencePort.DeploymentRepository {
	return &persistence.DeploymentRepositoryImpl{
		Db: dbService.GetDB(),
	}
}

func NewUserService() *services.UserService {
	return &services.UserService{}
}

func NewProjectService(projectRepo persistencePort.ProjectRepository, trainingDatasetRepo persistencePort.TrainingDatasetRepository, finetuneRepo persistencePort.FinetuneRepository, deploymentRepo persistencePort.DeploymentRepository) *services.ProjectService {
	return &services.ProjectService{
		ProjectRepository:         projectRepo,
		TrainingDatasetRepository: trainingDatasetRepo,
		FinetuneRepository:        finetuneRepo,
		DeploymentRepository:      deploymentRepo,
	}
}

func NewTrainingDatasetService() *services.TrainingDatasetService {
	return &services.TrainingDatasetService{}
}

func NewFinetuneService() *services.FinetuneService {
	return &services.FinetuneService{}
}

func NewFinetuneCompletionService(
	finetuneRepo persistencePort.FinetuneRepository,
	projectRepo persistencePort.ProjectRepository,
	ollamaLLMClient clientsPort.OllamaLLMClient,
) *services.FinetuneCompletionService {
	return services.NewFinetuneCompletionService(finetuneRepo, projectRepo, ollamaLLMClient)
}

func NewPromptAnalysisService(ollamaLLMClient clientsPort.OllamaLLMClient) *services.PromptAnalysisService {
	return &services.PromptAnalysisService{
		OllamaLLMClient: ollamaLLMClient,
	}
}

func NewDeploymentService(deploymentRepo persistencePort.DeploymentRepository, projectRepo persistencePort.ProjectRepository, finetuneRepo persistencePort.FinetuneRepository) *services.DeploymentService {
	return &services.DeploymentService{
		DeploymentRepository: deploymentRepo,
		ProjectRepository:    projectRepo,
		FinetuneRepository:   finetuneRepo,
	}
}

func NewJWTService() *services.JWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this-in-production"
	}
	return &services.JWTService{
		SecretKey: []byte(secretKey),
	}
}

func NewLoginUseCase(userRepo persistencePort.UserRepository, userService *services.UserService, jwtService *services.JWTService) in.LoginUseCase {
	return &use_cases.LoginUseCaseImpl{
		UserRepository: userRepo,
		UserService:    userService,
		JwtService:     jwtService,
	}
}

func NewCreateProjectUseCase(projectRepo persistencePort.ProjectRepository, projectService *services.ProjectService) in.CreateProjectUseCase {
	return &use_cases.CreateProjectUseCaseImpl{
		ProjectRepository: projectRepo,
		ProjectService:    projectService,
	}
}

func NewGetProjectUseCase(projectService *services.ProjectService) in.GetProjectUseCase {
	return &use_cases.GetProjectUseCaseImpl{
		ProjectService: projectService,
	}
}

func NewListProjectsUseCase(projectService *services.ProjectService) in.ListProjectsUseCase {
	return &use_cases.ListProjectsUseCaseImpl{
		ProjectService: projectService,
	}
}

func NewCreateTrainingDatasetUseCase(
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
	projectRepo persistencePort.ProjectRepository,
	corpusRepo persistencePort.CorpusRepository,
	promptRepo persistencePort.PromptRepository,
	trainingDatasetService *services.TrainingDatasetService,
	trainingDatasetJobClient clientsPort.TrainingDatasetJobClient,
) in.CreateTrainingDatasetUseCase {
	return &use_cases.CreateTrainingDatasetUseCaseImpl{
		TrainingDatasetRepository: trainingDatasetRepo,
		ProjectRepository:         projectRepo,
		CorpusRepository:          corpusRepo,
		PromptRepository:          promptRepo,
		TrainingDatasetService:    trainingDatasetService,
		TrainingDatasetJobClient:  trainingDatasetJobClient,
	}
}

func NewCreateFinetuneUseCase(
	finetuneRepo persistencePort.FinetuneRepository,
	projectRepo persistencePort.ProjectRepository,
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
	corpusRepo persistencePort.CorpusRepository,
	finetuneService *services.FinetuneService,
	trainingDatasetService *services.TrainingDatasetService,
	finetuneJobClient clientsPort.FinetuneJobClient,
	runpodClient clientsPort.RunpodClient,
) in.CreateFinetuneUseCase {
	return &use_cases.CreateFinetuneUseCaseImpl{
		FinetuneRepository:        finetuneRepo,
		ProjectRepository:         projectRepo,
		TrainingDatasetRepository: trainingDatasetRepo,
		CorpusRepository:          corpusRepo,
		FinetuneService:           finetuneService,
		TrainingDatasetService:    trainingDatasetService,
		FinetuneJobClient:         finetuneJobClient,
		RunpodClient:              runpodClient,
	}
}

func NewLoginController(loginUseCase in.LoginUseCase) *web.LoginController {
	return &web.LoginController{
		LoginUseCase: loginUseCase,
	}
}

func NewCreateProjectController(createProjectUseCase in.CreateProjectUseCase) *web.CreateProjectController {
	return &web.CreateProjectController{
		CreateProjectUseCase: createProjectUseCase,
	}
}

func NewGetProjectController(getProjectUseCase in.GetProjectUseCase) *web.GetProjectController {
	return &web.GetProjectController{
		GetProjectUseCase: getProjectUseCase,
	}
}

func NewListProjectsController(listProjectsUseCase in.ListProjectsUseCase) *web.ListProjectsController {
	return &web.ListProjectsController{
		ListProjectsUseCase: listProjectsUseCase,
	}
}

func NewCreateTrainingDatasetController(
	createTrainingDatasetUseCase in.CreateTrainingDatasetUseCase,
) *web.CreateTrainingDatasetController {
	return &web.CreateTrainingDatasetController{
		CreateTrainingDatasetUseCase: createTrainingDatasetUseCase,
	}
}

func NewUpdateTrainingDatasetStatusUseCase(
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
	trainingDatasetResultsClient clientsPort.TrainingDatasetResultsClient,
) in.UpdateTrainingDatasetStatusUseCase {
	return &use_cases.UpdateTrainingDatasetStatusUseCaseImpl{
		TrainingDatasetRepository:     trainingDatasetRepo,
		TrainingDatasetResultsClient: trainingDatasetResultsClient,
	}
}

func NewGetTrainingDatasetUseCase(
	trainingDatasetService *services.TrainingDatasetService,
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
	promptRepo persistencePort.PromptRepository,
	corpusRepo persistencePort.CorpusRepository,
) in.GetTrainingDatasetUseCase {
	return &use_cases.GetTrainingDatasetUseCaseImpl{
		TrainingDatasetService:    trainingDatasetService,
		TrainingDatasetRepository: trainingDatasetRepo,
		PromptRepository:          promptRepo,
		CorpusRepository:          corpusRepo,
	}
}

func NewUpdateTrainingDatasetStatusController(updateTrainingDatasetStatusUseCase in.UpdateTrainingDatasetStatusUseCase) *web.UpdateTrainingDatasetStatusController {
	return &web.UpdateTrainingDatasetStatusController{
		UpdateTrainingDatasetStatusUseCase: updateTrainingDatasetStatusUseCase,
	}
}

func NewUpdateFinetuneStatusUseCase(
	finetuneRepo persistencePort.FinetuneRepository,
) in.UpdateFinetuneStatusUseCase {
	return &use_cases.UpdateFinetuneStatusUseCaseImpl{
		FinetuneRepository: finetuneRepo,
	}
}

func NewUpdateFinetuneStatusController(updateFinetuneStatusUseCase in.UpdateFinetuneStatusUseCase) *web.UpdateFinetuneStatusController {
	return &web.UpdateFinetuneStatusController{
		UpdateFinetuneStatusUseCase: updateFinetuneStatusUseCase,
	}
}

func NewGetFinetuneUseCase(finetuneRepo persistencePort.FinetuneRepository) in.GetFinetuneUseCase {
	return &use_cases.GetFinetuneUseCaseImpl{
		FinetuneRepository: finetuneRepo,
	}
}

func NewFinetuneCompletionUseCase(
	finetuneCompletionService *services.FinetuneCompletionService,
) in.FinetuneCompletionUseCase {
	return use_cases.NewFinetuneCompletionUseCaseImpl(finetuneCompletionService)
}

func NewGetFinetuneController(getFinetuneUseCase in.GetFinetuneUseCase) *web.GetFinetuneController {
	return &web.GetFinetuneController{
		GetFinetuneUseCase: getFinetuneUseCase,
	}
}

func NewDownloadModelUseCase(finetuneRepo persistencePort.FinetuneRepository, downloadModelClient clientsPort.DownloadModelClient) in.DownloadModelUseCase {
	return &use_cases.DownloadModelUseCaseImpl{
		FinetuneRepository:  finetuneRepo,
		DownloadModelClient: downloadModelClient,
	}
}

func NewAnalyzePromptUseCase(promptAnalysisService *services.PromptAnalysisService) in.AnalyzePromptUseCase {
	return &use_cases.AnalyzePromptUseCaseImpl{
		PromptAnalysisService: promptAnalysisService,
	}
}

func NewDownloadModelController(downloadModelUseCase in.DownloadModelUseCase) *web.DownloadModelController {
	return &web.DownloadModelController{
		DownloadModelUseCase: downloadModelUseCase,
	}
}

func NewAnalyzePromptController(analyzePromptUseCase in.AnalyzePromptUseCase) *web.AnalyzePromptController {
	return &web.AnalyzePromptController{
		AnalyzePromptUseCase: analyzePromptUseCase,
	}
}

func NewGetTrainingDatasetController(getTrainingDatasetUseCase in.GetTrainingDatasetUseCase) *web.GetTrainingDatasetController {
	return &web.GetTrainingDatasetController{
		GetTrainingDatasetUseCase: getTrainingDatasetUseCase,
	}
}

func NewDownloadTrainingDatasetUseCase(
	trainingDatasetService *services.TrainingDatasetService,
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
	projectRepo persistencePort.ProjectRepository,
) in.DownloadTrainingDatasetUseCase {
	return &use_cases.DownloadTrainingDatasetUseCaseImpl{
		TrainingDatasetService:    trainingDatasetService,
		TrainingDatasetRepository: trainingDatasetRepo,
		ProjectRepository:         projectRepo,
	}
}

func NewDownloadTrainingDatasetController(downloadTrainingDatasetUseCase in.DownloadTrainingDatasetUseCase) *web.DownloadTrainingDatasetController {
	return &web.DownloadTrainingDatasetController{
		DownloadTrainingDatasetUseCase: downloadTrainingDatasetUseCase,
	}
}

func NewUploadTrainingDatasetUseCase(
	trainingDatasetService *services.TrainingDatasetService,
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
) in.UploadTrainingDatasetUseCase {
	return &use_cases.UploadTrainingDatasetUseCaseImpl{
		TrainingDatasetService:    trainingDatasetService,
		TrainingDatasetRepository: trainingDatasetRepo,
	}
}

func NewUploadTrainingDatasetController(uploadTrainingDatasetUseCase in.UploadTrainingDatasetUseCase) *web.UploadTrainingDatasetController {
	return &web.UploadTrainingDatasetController{
		UploadTrainingDatasetUseCase: uploadTrainingDatasetUseCase,
	}
}

func NewUploadNewTrainingDatasetVersionUseCase(
	trainingDatasetService *services.TrainingDatasetService,
	trainingDatasetRepo persistencePort.TrainingDatasetRepository,
) in.UploadNewTrainingDatasetVersionUseCase {
	return &use_cases.UploadNewTrainingDatasetVersionUseCaseImpl{
		TrainingDatasetService:    trainingDatasetService,
		TrainingDatasetRepository: trainingDatasetRepo,
	}
}

func NewUploadNewTrainingDatasetVersionController(uploadNewTrainingDatasetVersionUseCase in.UploadNewTrainingDatasetVersionUseCase) *web.UploadNewTrainingDatasetVersionController {
	return &web.UploadNewTrainingDatasetVersionController{
		UploadNewTrainingDatasetVersionUseCase: uploadNewTrainingDatasetVersionUseCase,
	}
}

func NewCreateFinetuneController(createFinetuneUseCase in.CreateFinetuneUseCase) *web.CreateFinetuneController {
	return &web.CreateFinetuneController{
		CreateFinetuneUseCase: createFinetuneUseCase,
	}
}

func NewFinetuneCompletionController(finetuneCompletionUseCase in.FinetuneCompletionUseCase) *web.FinetuneCompletionController {
	return &web.FinetuneCompletionController{
		FinetuneCompletionUseCase: finetuneCompletionUseCase,
	}
}

func NewCreateDeploymentUseCase(deploymentRepo persistencePort.DeploymentRepository, deploymentService *services.DeploymentService) in.CreateDeploymentUseCase {
	return &use_cases.CreateDeploymentUseCaseImpl{
		DeploymentRepository: deploymentRepo,
		DeploymentService:    deploymentService,
	}
}

func NewCreateDeploymentController(createDeploymentUseCase in.CreateDeploymentUseCase) *web.CreateDeploymentController {
	return &web.CreateDeploymentController{
		CreateDeploymentUseCase: createDeploymentUseCase,
	}
}

func NewGetDeploymentUseCase(deploymentRepo persistencePort.DeploymentRepository, deploymentService *services.DeploymentService) in.GetDeploymentUseCase {
	return &use_cases.GetDeploymentUseCaseImpl{
		DeploymentRepository: deploymentRepo,
		DeploymentService:    deploymentService,
	}
}

func NewGetDeploymentController(getDeploymentUseCase in.GetDeploymentUseCase) *web.GetDeploymentController {
	return &web.GetDeploymentController{
		GetDeploymentUseCase: getDeploymentUseCase,
	}
}

func NewExternalAPIMiddleware() *server.ExternalAPIMiddleware {
	return &server.ExternalAPIMiddleware{}
}

func NewTrainingDatasetJobClient() clientsPort.TrainingDatasetJobClient {
	client, err := clients.NewTrainingDatasetJobClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewTrainingDatasetResultsClient() clientsPort.TrainingDatasetResultsClient {
	client, err := clients.NewTrainingDatasetResultsClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewFinetuneJobClient() clientsPort.FinetuneJobClient {
	client, err := clients.NewFinetuneJobClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewRunpodClient() clientsPort.RunpodClient {
	client, err := clients.NewRunpodClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewDownloadModelClient() clientsPort.DownloadModelClient {
	client, err := clients.NewDownloadModelClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewOllamaLLMClient() clientsPort.OllamaLLMClient {
	client, err := clients.NewOllamaLLMClientImpl()
	if err != nil {
		panic(err)
	}
	return client
}

func NewAuthMiddleware(jwtService *services.JWTService) *server.AuthMiddleware {
	return &server.AuthMiddleware{
		JwtService: jwtService,
	}
}

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewProjectRepository),
	fx.Provide(NewTrainingDatasetRepository),
	fx.Provide(NewCorpusRepository),
	fx.Provide(NewPromptRepository),
	fx.Provide(NewFinetuneRepository),
	fx.Provide(NewDeploymentRepository),
	fx.Provide(NewTrainingDatasetJobClient),
	fx.Provide(NewTrainingDatasetResultsClient),
	fx.Provide(NewFinetuneJobClient),
	fx.Provide(NewRunpodClient),
	fx.Provide(NewDownloadModelClient),
	fx.Provide(NewOllamaLLMClient),
	fx.Provide(NewUserService),
	fx.Provide(NewProjectService),
	fx.Provide(NewTrainingDatasetService),
	fx.Provide(NewFinetuneService),
	fx.Provide(NewFinetuneCompletionService),
	fx.Provide(NewPromptAnalysisService),
	fx.Provide(NewDeploymentService),
	fx.Provide(NewJWTService),
	fx.Provide(NewLoginUseCase),
	fx.Provide(NewCreateProjectUseCase),
	fx.Provide(NewGetProjectUseCase),
	fx.Provide(NewListProjectsUseCase),
	fx.Provide(NewCreateTrainingDatasetUseCase),
	fx.Provide(NewCreateFinetuneUseCase),
	fx.Provide(NewGetTrainingDatasetUseCase),
	fx.Provide(NewDownloadTrainingDatasetUseCase),
	fx.Provide(NewUploadTrainingDatasetUseCase),
	fx.Provide(NewUploadNewTrainingDatasetVersionUseCase),
	fx.Provide(NewUpdateTrainingDatasetStatusUseCase),
	fx.Provide(NewUpdateFinetuneStatusUseCase),
	fx.Provide(NewGetFinetuneUseCase),
	fx.Provide(NewFinetuneCompletionUseCase),
	fx.Provide(NewDownloadModelUseCase),
	fx.Provide(NewAnalyzePromptUseCase),
	fx.Provide(NewCreateDeploymentUseCase),
	fx.Provide(NewGetDeploymentUseCase),
	fx.Provide(NewLoginController),
	fx.Provide(NewCreateProjectController),
	fx.Provide(NewGetProjectController),
	fx.Provide(NewListProjectsController),
	fx.Provide(NewCreateTrainingDatasetController),
	fx.Provide(NewCreateFinetuneController),
	fx.Provide(NewGetTrainingDatasetController),
	fx.Provide(NewDownloadTrainingDatasetController),
	fx.Provide(NewUploadTrainingDatasetController),
	fx.Provide(NewUploadNewTrainingDatasetVersionController),
	fx.Provide(NewUpdateTrainingDatasetStatusController),
	fx.Provide(NewUpdateFinetuneStatusController),
	fx.Provide(NewGetFinetuneController),
	fx.Provide(NewFinetuneCompletionController),
	fx.Provide(NewDownloadModelController),
	fx.Provide(NewAnalyzePromptController),
	fx.Provide(NewCreateDeploymentController),
	fx.Provide(NewGetDeploymentController),
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewExternalAPIMiddleware),
)