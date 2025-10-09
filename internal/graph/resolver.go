package graph

import (
	"tsv-golang/internal/repository"
	"tsv-golang/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repositories *repository.Repositories
	Service      *service.Service
}
