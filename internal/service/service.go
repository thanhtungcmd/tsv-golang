package service

import (
	"tsv-golang/internal/repository"
)

type Service struct {
	repo        *repository.Repositories
	UserService UserServiceInterface
}

func NewService(repo *repository.Repositories) *Service {
	return &Service{
		repo:        repo,
		UserService: UserServiceInit(repo),
	}
}
