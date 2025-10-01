package env

import (
	"os"
)

type env struct {
}

func Env() *env {
	return &env{}
}

func (h *env) IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func (h *env) IsLocal() bool {
	return os.Getenv("APP_ENV") == "local"
}

func (h *env) IsTest() bool {
	return os.Getenv("APP_ENV") == "test"
}

func (h *env) IsModeDebug() bool {
	return os.Getenv("GIN_MODE") == "debug"
}
