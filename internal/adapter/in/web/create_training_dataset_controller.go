package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type CreateTrainingDatasetController struct {
	CreateTrainingDatasetUseCase in.CreateTrainingDatasetUseCase
}


func (c *CreateTrainingDatasetController) CreateTrainingDataset(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	projectIDParam := ctx.Param("project_id")
	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID format",
		})
		return
	}

	var request CreateTrainingDatasetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set default values for generate model and runner if not provided
	generateModel := request.GenerateModel
	if generateModel == "" {
		generateModel = "qwen3:30b-a3b-instruct-2507-q4_K_M"
	}
	generateModelRunner := request.GenerateModelRunner
	if generateModelRunner == "" {
		generateModelRunner = "runpod_ollama"
	}

	command := in.CreateTrainingDatasetCommand{
		UserID:                 userID,
		ProjectID:              projectID,
		CorpusName:             request.CorpusName,
		InputField:             request.InputField,
		OutputField:            request.OutputField,
		LanguageISO:            request.LanguageISO,
		FieldNames:             request.FieldNames,
		GeneratePrompt:         request.GeneratePrompt,
		GenerateExamplesNumber: request.GenerateExamplesNumber,
		GenerateModel:          generateModel,
		GenerateModelRunner:    generateModelRunner,
	}

	result, err := c.CreateTrainingDatasetUseCase.Execute(ctx.Request.Context(), command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := ToCreateTrainingDatasetResponse(result)
	ctx.JSON(http.StatusCreated, response)
}