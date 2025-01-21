package util

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTUtil struct {}

func (j *JWTUtil) GetToken(ctx context.Context, claims map[string]interface{}) (*string, error) {
	var secret = []byte(os.Getenv("JWT_SECRET_KEY")) 
	builder := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(12*time.Hour))

	for k, v := range claims {
		builder = builder.Claim(k, v)
	}
	token, err := builder.Build()
  if err != nil {
    return nil, dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("error building token: %s", err.Error()),
		}
  }

  tokenByte, err := jwt.NewSerializer().Sign(jwt.WithKey(jwa.HS256, secret)).Serialize(token)
  if err != nil {
    return nil, dto.ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: fmt.Sprintf("error signing and serializing token: %s", err.Error()),
		}
  }

	tokenString := string(tokenByte)
	return &tokenString, nil
}

func (j *JWTUtil) ValidateToken(ctx context.Context, tokenString string) (jwt.Token, error) {
	var secret = []byte(os.Getenv("JWT_SECRET_KEY")) 
	token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256, secret))
  if err != nil {
    return nil, dto.ErrorResponse{
			Status: http.StatusUnauthorized,
			Message: fmt.Sprintf("invalid token: %s", err),
		}
  }
	return token, nil
}