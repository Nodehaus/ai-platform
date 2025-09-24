package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ai-platform/internal/application/port/in"
)

type GetFinetuneController struct {
	GetFinetuneUseCase in.GetFinetuneUseCase
}

func (c *GetFinetuneController) GetFinetune(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	projectIDStr := ctx.Param("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID format",
		})
		return
	}

	finetuneIDStr := ctx.Param("finetune_id")
	finetuneID, err := uuid.Parse(finetuneIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid finetune ID format",
		})
		return
	}

	command := in.GetFinetuneCommand{
		ProjectID:  projectID,
		FinetuneID: finetuneID,
		OwnerID:    userID,
	}

	result, err := c.GetFinetuneUseCase.GetFinetune(ctx.Request.Context(), command)
	if err != nil {
		if err.Error() == "finetune not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Finetune not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch finetune",
		})
		return
	}

	if result == nil || result.Finetune == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Finetune not found",
		})
		return
	}

	response := ToGetFinetuneResponse(result.Finetune)
	ctx.JSON(http.StatusOK, response)
}