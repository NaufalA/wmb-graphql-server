package service

import (
	"context"
	"net/http"

	"github.com/NaufalA/wmb-graphql-server/graph/model"
	"github.com/NaufalA/wmb-graphql-server/internal/collection"
	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/NaufalA/wmb-graphql-server/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	GetUser(context.Context, dto.GetUserRequest) (*collection.User, error)
	CreateUser(ctx context.Context, input model.CreateUserInput) (*collection.User, error)
	UpdateUser(context.Context, collection.User) (*collection.User, error)
}

type AuthService struct {
	logger         *logrus.Entry
	userRepository UserRepository
}

func NewAuthService(
	logrus *logrus.Logger,
	userRepository UserRepository,
) *AuthService {
	return &AuthService{
		logger:         logrus.WithField("location", "AuthService"),
		userRepository: userRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, _ := s.userRepository.GetUser(ctx, dto.GetUserRequest{Email: request.Email})
	if user == nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "wrong email or password",
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}
	passwordUtil := util.PasswordUtil{}
	valid := passwordUtil.CheckPasswordHash(request.Password, *user.PasswordHash)
	if !valid {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "wrong email or password",
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	jwtUtil := util.JWTUtil{}
	claims := make(map[string]interface{})
	claims["user"] = model.User{
		ID: user.ID.Hex(),
		Email: user.Email,
		FullName: user.FullName,
		Role: user.Role,
	}
	tokenString, err := jwtUtil.GetToken(ctx, claims)
	if err != nil  {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	return &dto.LoginResponse{
		Token: *tokenString,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, request model.CreateUserInput) (*model.User, error) {
	exist, _ := s.userRepository.GetUser(ctx, dto.GetUserRequest{Email: request.Email})
	if exist != nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusConflict,
			Message: "email already exist",
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	response, err := s.userRepository.CreateUser(ctx, request)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	return &model.User{
		ID:         response.ID.Hex(),
		Email:      response.Email,
		FullName:   response.FullName,
		Role:       response.Role,
		CreateTime: response.CreateTime,
		UpdateTime: response.UpdateTime,
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	user, _ := s.userRepository.GetUser(ctx, dto.GetUserRequest{Email: request.Email})
	if user == nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "user doesn't exist",
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	passwordUtil := util.PasswordUtil{}
	newPassword := bson.NewObjectID().Hex()
	newPasswordHash, err := passwordUtil.HashPassword(newPassword)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}
	user.PasswordHash = &newPasswordHash
	_, err = s.userRepository.UpdateUser(ctx, *user)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		s.logger.Error(errorResponse.Message)
		return nil, errorResponse
	}

	return &dto.ResetPasswordResponse{
		NewPassword: newPassword,
	}, nil
}
