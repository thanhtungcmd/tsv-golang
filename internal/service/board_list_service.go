package service

import (
	"errors"
	"fmt"
	"tsv-golang/internal/graph/model"
	"tsv-golang/internal/repository"
	"tsv-golang/pkg/datetime"

	"github.com/mitchellh/mapstructure"
)

type BoardListService struct {
	repo *repository.Repositories
}

type BoardListServiceInterface interface {
	CreateBoardList(userLogin string, input *model.BoardListInput) (*model.BoardList, error)
	UpdateBoardList(userLogin string, id string, input model.BoardListUpdateInput) (*model.BoardList, error)
	DeleteBoardList(userLogin string, id string) (*string, error)
	ListBoardList(request *model.ListBoardListRequest) ([]*model.BoardList, error)
	GetBoardListByID(id string) (*model.BoardList, error)
}

func BoardListServiceInit(repo *repository.Repositories) *BoardListService {
	return &BoardListService{
		repo: repo,
	}
}

var _ BoardListServiceInterface = &BoardListService{}

func (r *BoardListService) CreateBoardList(userLogin string, input *model.BoardListInput) (*model.BoardList, error) {
	board := r.repo.Board.FindById(input.BoardID)
	if board == nil {
		return nil, errors.New(fmt.Sprintf("board not found"))
	}
	result := &model.BoardList{}
	err := mapstructure.Decode(input, &result)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	result.UseYn = &STATUS_ACTIVE
	result.CreatedAt = &timeNow
	result.UpdatedAt = &timeNow
	result.CreatedBy = &userLogin
	result.UpdatedBy = &userLogin
	result, err = r.repo.BoardList.CreateAndReturn(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *BoardListService) UpdateBoardList(userLogin string, id string, input model.BoardListUpdateInput) (*model.BoardList, error) {
	data := r.repo.BoardList.FindById(id)
	err := mapstructure.Decode(input, &data)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	data.UpdatedAt = &timeNow
	data.UpdatedBy = &userLogin
	err = r.repo.BoardList.UpdateByConditions(id, *data, []string{"name", "sort_order"}...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *BoardListService) DeleteBoardList(userLogin string, id string) (*string, error) {
	data := r.repo.BoardList.FindById(id)
	if data == nil {
		return nil, fmt.Errorf("board list not found")
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	data.UseYn = &STATUS_IN_ACTIVE
	data.UpdatedAt = &timeNow
	data.UpdatedBy = &userLogin
	err := r.repo.BoardList.UpdateByConditions(id, *data, []string{"use_yn"}...)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *BoardListService) ListBoardList(request *model.ListBoardListRequest) ([]*model.BoardList, error) {
	result := r.repo.BoardList.GetList(request)
	return result, nil
}

func (r *BoardListService) GetBoardListByID(id string) (*model.BoardList, error) {
	result := r.repo.BoardList.FindById(id)
	return result, nil
}
