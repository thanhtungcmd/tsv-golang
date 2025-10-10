package service

import (
	"tsv-golang/internal/repository"
)

type Service struct {
	repo             *repository.Repositories
	UserService      UserServiceInterface
	BoardService     BoardServiceInterface
	BoardListService BoardListServiceInterface
	BoardCardService BoardCardServiceInterface
}

func NewService(repo *repository.Repositories) *Service {
	return &Service{
		repo:             repo,
		UserService:      UserServiceInit(repo),
		BoardService:     BoardServiceInit(repo),
		BoardListService: BoardListServiceInit(repo),
		BoardCardService: BoardCardServiceInit(repo),
	}
}
