package service

import "tsv-golang/internal/repository"

type BoardListService struct {
	repo *repository.Repositories
}

type BoardListServiceInterface interface {
}

func BoardListServiceInit(repo *repository.Repositories) *BoardListService {
	return &BoardListService{
		repo: repo,
	}
}

var _ BoardListServiceInterface = &BoardListService{}
