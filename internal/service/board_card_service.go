package service

import "tsv-golang/internal/repository"

type BoardCardService struct {
	repo *repository.Repositories
}

type BoardCardServiceInterface interface {
}

func BoardCardServiceInit(repo *repository.Repositories) *BoardCardService {
	return &BoardCardService{
		repo: repo,
	}
}

var _ BoardCardServiceInterface = &BoardCardService{}
