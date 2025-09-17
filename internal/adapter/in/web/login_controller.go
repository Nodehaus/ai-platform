package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ai-platform/internal/application/port/in"
)

type LoginController struct {
	loginUseCase in.LoginUseCase
}

func NewLoginController(loginUseCase in.LoginUseCase) *LoginController {
	return &LoginController{
		loginUseCase: loginUseCase,
	}
}

func (c *LoginController) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	command := in.LoginCommand{
		Email:    request.Email,
		Password: request.Password,
	}

	loginResult, err := c.loginUseCase.Login(command)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := NewLoginResponse(loginResult.User, loginResult.Token, "Login successful")
	ctx.JSON(http.StatusOK, response)
}