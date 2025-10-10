package direction

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
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

	claims, err := tokenValid(token, apiSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %s", err)
	}

	username, ok := claims["username"].(string)
	if ok {
		ctx = context.WithValue(ctx, "userLogin", username)
	}

	return next(ctx)
}

func AuthContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		ctx := context.WithValue(c.Request.Context(), "token", tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func tokenValid(tokenString string, secretKey string) (jwt.MapClaims, error) {
	tokens, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	alg, ok := tokens.Header["alg"].(string)
	if !ok || alg != "HS256" {
		return nil, AuthInvalidAlg
	}
	typ, ok := tokens.Header["typ"].(string)
	if !ok || typ != "JWT" {
		return nil, AuthInvalidTyp
	}
	claims, ok := tokens.Claims.(jwt.MapClaims)
	if !ok || !tokens.Valid {
		return nil, AuthInvalidExp
	}
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, AuthInvalidExp
	}
	return claims, nil
}
