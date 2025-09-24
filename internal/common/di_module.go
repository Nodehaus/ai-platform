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

func NewUserService() *services.UserService {
	return &services.UserService{}
}

func NewProjectService(projectRepo persistencePort.ProjectRepository, trainingDatasetRepo persistencePort.TrainingDatasetRepository) *services.ProjectService {
	return &services.ProjectService{
		ProjectRepository:         projectRepo,
		TrainingDatasetRepository: trainingDatasetRepo,
	}
}

func NewTrainingDatasetService() *services.TrainingDatasetService {
	return &services.TrainingDatasetService{}
}

func NewFinetuneService() *services.FinetuneService {
	return &services.FinetuneService{}
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
	finetuneService *services.FinetuneService,
) in.CreateFinetuneUseCase {
	return &use_cases.CreateFinetuneUseCaseImpl{
		FinetuneRepository:        finetuneRepo,
		ProjectRepository:         projectRepo,
		TrainingDatasetRepository: trainingDatasetRepo,
		FinetuneService:           finetuneService,
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

func NewGetTrainingDatasetController(getTrainingDatasetUseCase in.GetTrainingDatasetUseCase) *web.GetTrainingDatasetController {
	return &web.GetTrainingDatasetController{
		GetTrainingDatasetUseCase: getTrainingDatasetUseCase,
	}
}

func NewCreateFinetuneController(createFinetuneUseCase in.CreateFinetuneUseCase) *web.CreateFinetuneController {
	return &web.CreateFinetuneController{
		CreateFinetuneUseCase: createFinetuneUseCase,
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
	fx.Provide(NewTrainingDatasetJobClient),
	fx.Provide(NewTrainingDatasetResultsClient),
	fx.Provide(NewUserService),
	fx.Provide(NewProjectService),
	fx.Provide(NewTrainingDatasetService),
	fx.Provide(NewFinetuneService),
	fx.Provide(NewJWTService),
	fx.Provide(NewLoginUseCase),
	fx.Provide(NewCreateProjectUseCase),
	fx.Provide(NewGetProjectUseCase),
	fx.Provide(NewListProjectsUseCase),
	fx.Provide(NewCreateTrainingDatasetUseCase),
	fx.Provide(NewCreateFinetuneUseCase),
	fx.Provide(NewGetTrainingDatasetUseCase),
	fx.Provide(NewUpdateTrainingDatasetStatusUseCase),
	fx.Provide(NewLoginController),
	fx.Provide(NewCreateProjectController),
	fx.Provide(NewGetProjectController),
	fx.Provide(NewListProjectsController),
	fx.Provide(NewCreateTrainingDatasetController),
	fx.Provide(NewCreateFinetuneController),
	fx.Provide(NewGetTrainingDatasetController),
	fx.Provide(NewUpdateTrainingDatasetStatusController),
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewExternalAPIMiddleware),
)