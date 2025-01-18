package controller

import (
	"context"
	"net/http"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type authService interface {
	Login(context.Context, dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, model.CreateUserInput) (*model.User, error)
	ResetPassword(context.Context, dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

type AuthController struct {
	logger *logrus.Entry
	authService authService
}

func NewAuthController(
	logrus *logrus.Logger,
	authService authService,
	) *AuthController {
	return &AuthController{
		logger: logrus.WithField("location", "AuthController"),
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	request := dto.LoginRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		c.logger.Error(err)
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, dto.ErrorResponse{
			Status: status,
			Message: err.Error(),
		})
		return
	}

	response, err := c.authService.Login(ctx, request)
	if err != nil {
		errorResponse := err.(dto.ErrorResponse)
		c.logger.Error(errorResponse)
		ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) Register(ctx *gin.Context) {
	request := model.CreateUserInput{}
	err := ctx.BindJSON(&request)
	if err != nil {
		c.logger.Error(err)
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, dto.ErrorResponse{
			Status: status,
			Message: err.Error(),
		})
		return
	}

	response, err := c.authService.Register(ctx, request)
	if err != nil {
		errorResponse := err.(dto.ErrorResponse)
		c.logger.Error(errorResponse)
		ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	request := dto.ResetPasswordRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		c.logger.Error(err)
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, dto.ErrorResponse{
			Status: status,
			Message: err.Error(),
		})
		return
	}

	response, err := c.authService.ResetPassword(ctx, request)
	if err != nil {
		errorResponse := err.(dto.ErrorResponse)
		c.logger.Error(errorResponse)
		ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
		return
	}

	ctx.JSON(http.StatusOK, response)
}