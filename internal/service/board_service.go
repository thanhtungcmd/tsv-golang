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
	UpdateBoard(userLogin string, id string, input model.BoardUpdateInput) (*model.Board, error)
	//DeleteBoard(userLogin string, id string) error
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
	result, err = r.repo.Board.CreateAndReturn(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *BoardService) UpdateBoard(userLogin string, id string, input model.BoardUpdateInput) (*model.Board, error) {
	data := r.repo.Board.FindById(id)
	err := mapstructure.Decode(input, &data)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	data.UpdatedAt = &timeNow
	data.UpdatedBy = &userLogin
	err = r.repo.Board.UpdateByConditions(id, *data, []string{"name", "description", "sort_order"}...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//func (r *BoardService) DeleteBoard(userLogin string, id string) error {
//	data := r.repo.Board.FindById(id)
//	if err != nil {
//		return nil, err
//	}
//	timeNow := datetime.Datetime().TimeNow().ToString()
//	data.UpdatedAt = &timeNow
//	data.UpdatedBy = &userLogin
//	err = r.repo.Board.UpdateByConditions(id, *data, []string{"use_yn"}...)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
