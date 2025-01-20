package middleware

import (
	"net/http"
	"strings"

	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/NaufalA/wmb-graphql-server/pkg/util"
	"github.com/gin-gonic/gin"
)

func AuthTokenMiddleware() gin.HandlerFunc  {
	return func(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
		if authorization == "" {
			errorResponse := dto.ErrorResponse{
				Status: http.StatusUnauthorized,
				Message: "unauthorized. please login",
			}
			ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}
		authorizationSplit := strings.Split(authorization, " ")
		if len(authorizationSplit) < 2 {
			errorResponse := dto.ErrorResponse{
				Status: http.StatusUnauthorized,
				Message: "invalid token",
			}
			ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}
		tokenString := strings.Split(authorization, " ")[1]
		jwtUtil := util.JWTUtil{}
		token, err := jwtUtil.ValidateToken(ctx, tokenString)
		if err != nil {
			errorResponse := err.(dto.ErrorResponse)
			ctx.AbortWithStatusJSON(errorResponse.Status, errorResponse)
			return
		}
		userData, exist := token.Get("user")
		if exist {
			ctx.Set("user", userData)
		}
		ctx.Next()
	}
}