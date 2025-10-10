package direction

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt/v4"
)

var (
	AuthInvalidAlg = errors.New("algorithm is not valid")
	AuthInvalidTyp = errors.New("typ is not valid")
	AuthInvalidExp = errors.New("expire time is not valid")
)

func Authen(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	token, _ := ctx.Value("token").(string)
	if token == "" {
		return nil, fmt.Errorf("unauthorized: missing token")
	}

	apiSecretKey := os.Getenv("SECRET_KEY")

	err = tokenValid(token, apiSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %s", err)
	}

	return next(ctx)
}

func tokenValid(tokenString string, secretKey string) error {
	tokens, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}

	alg, ok := tokens.Header["alg"].(string)
	if !ok || alg != "HS256" {
		return AuthInvalidAlg
	}
	typ, ok := tokens.Header["typ"].(string)
	if !ok || typ != "JWT" {
		return AuthInvalidTyp
	}
	claims, ok := tokens.Claims.(jwt.MapClaims)
	if !ok || !tokens.Valid {
		return AuthInvalidExp
	}
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return AuthInvalidExp
	}
	return nil
}
