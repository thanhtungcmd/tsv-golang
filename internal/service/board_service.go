package service

import (
	"tsv-golang/internal/graph/model"
	"tsv-golang/internal/repository"
	"tsv-golang/pkg/datetime"

	"github.com/mitchellh/mapstructure"
)

type BoardService struct {
	repo *repository.Repositories
}

type BoardServiceInterface interface {
	CreateBoard(userLogin string, input *model.BoardInput) (*model.Board, error)
}

func BoardServiceInit(repo *repository.Repositories) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

var _ BoardServiceInterface = &BoardService{}

func (r *BoardService) CreateBoard(userLogin string, input *model.BoardInput) (*model.Board, error) {
	result := &model.Board{}
	err := mapstructure.Decode(input, &result)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	result.CreatedAt = &timeNow
	result.UpdatedAt = &timeNow
	result.CreatedBy = &userLogin
	result.UpdatedBy = &userLogin
	
	return result, nil
}
