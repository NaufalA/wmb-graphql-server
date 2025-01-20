package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NaufalA/wmb-graphql-server/internal/dto"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTUtil struct {}

var secret = []byte("your_secret_key") 

func (j *JWTUtil) GetToken(ctx context.Context, claims map[string]interface{}) (*string, error) {
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
	token, err := jwt.ParseString(tokenString, jwt.WithKey(jwa.HS256, secret))
  if err != nil {
    return nil, dto.ErrorResponse{
			Status: http.StatusUnauthorized,
			Message: fmt.Sprintf("invalid token: %s", err),
		}
  }
	return token, nil
}