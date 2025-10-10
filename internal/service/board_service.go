package service

import "tsv-golang/internal/repository"

type BoardService struct {
	repo *repository.Repositories
}

type BoardServiceInterface interface {
}

func BoardServiceInit(repo *repository.Repositories) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

var _ BoardServiceInterface = &BoardService{}
