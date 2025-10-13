package service

import (
	"tsv-golang/internal/repository"
)

type Service struct {
	repo             *repository.Repositories
	UserService      UserServiceInterface
	ProjectService   ProjectServiceInterface
	BoardService     BoardServiceInterface
	BoardListService BoardListServiceInterface
	BoardCardService BoardCardServiceInterface
}

var (
	STATUS_ACTIVE    = 1
	STATUS_IN_ACTIVE = 0
)

func NewService(repo *repository.Repositories) *Service {
	return &Service{
		repo:             repo,
		UserService:      UserServiceInit(repo),
		ProjectService:   ProjectServiceInit(repo),
		BoardService:     BoardServiceInit(repo),
		BoardListService: BoardListServiceInit(repo),
		BoardCardService: BoardCardServiceInit(repo),
	}
}
