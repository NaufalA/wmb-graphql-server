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
			status := http.StatusUnauthorized
			ctx.AbortWithStatusJSON(status, dto.ErrorResponse{
				Status: status,
				Message: "unauthorized. please login",
			})
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