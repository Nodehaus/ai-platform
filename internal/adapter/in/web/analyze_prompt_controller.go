package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-platform/internal/application/port/in"
)

type AnalyzePromptController struct {
	AnalyzePromptUseCase in.AnalyzePromptUseCase
}

func (c *AnalyzePromptController) AnalyzePrompt(ctx *gin.Context) {
	var request AnalyzePromptRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	command := in.AnalyzePromptCommand{
		Prompt: request.Prompt,
	}

	result, err := c.AnalyzePromptUseCase.AnalyzePrompt(command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := NewAnalyzePromptResponse(
		result.AnalysisResult,
		result.JSONObjectFields,
		result.InputField,
		result.OutputField,
		result.ExpectedOutputSizeChars,
	)
	ctx.JSON(http.StatusOK, response)
}
