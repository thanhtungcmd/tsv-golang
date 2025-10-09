package filter

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"tsv-golang/internal/route"
	"tsv-golang/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var skipURLs = []string{
	"/api/v1/auth/login",
	"/api/v1/",
}

var (
	AuthInvalidAlg = errors.New("algorithm is not valid")
	AuthInvalidTyp = errors.New("typ is not valid")
	AuthInvalidExp = errors.New("expire time is not valid")
)

func AuthFilter(c *gin.Context) {
	path := c.Request.URL.Path
	if isSkipURL(path) {
		c.Next()
		return
	}

	apiSecretKey := os.Getenv("SECRET_KEY")
	defer c.Request.Body.Close()

	if err := tokenValid(c, apiSecretKey); err != nil {
		log.Error(err)
		res := &route.ResponseStruct{
			ErrorCode: 401,
			Message:   "Unauthorized",
		}
		route.Response(c, res, http.StatusUnauthorized)
		c.Abort()
		return
	}
	c.Next()
}

func tokenValid(c *gin.Context, secretKey string) error {
	tokenString := extractToken(c)
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

func extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func isSkipURL(path string) bool {
	for _, url := range skipURLs {
		if strings.HasSuffix(url, "**") {
			prefix := strings.TrimSuffix(url, "**")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		} else {
			if path == url {
				return true
			}
		}
	}
	return false
}
