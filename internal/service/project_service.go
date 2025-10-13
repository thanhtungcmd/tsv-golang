package service

import (
	"fmt"
	"tsv-golang/internal/graph/model"
	"tsv-golang/internal/repository"
	"tsv-golang/pkg/datetime"

	"github.com/mitchellh/mapstructure"
)

type ProjectService struct {
	repo *repository.Repositories
}

type ProjectServiceInterface interface {
	CreateProject(userLogin string, input *model.ProjectInput) (*model.Project, error)
	UpdateProject(userLogin string, id string, input model.ProjectUpdateInput) (*model.Project, error)
	DeleteProject(userLogin string, id string) (*string, error)
	ListProject(request *model.ListProjectRequest) ([]*model.Project, error)
	GetProjectByID(id string) (*model.Project, error)
}

func ProjectServiceInit(repo *repository.Repositories) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

var _ ProjectServiceInterface = &ProjectService{}

func (r *ProjectService) CreateProject(userLogin string, input *model.ProjectInput) (*model.Project, error) {
	result := &model.Project{}
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
	result, err = r.repo.Project.CreateAndReturn(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProjectService) UpdateProject(userLogin string, id string, input model.ProjectUpdateInput) (*model.Project, error) {
	data := r.repo.Project.FindById(id)
	err := mapstructure.Decode(input, &data)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	data.UpdatedAt = &timeNow
	data.UpdatedBy = &userLogin
	err = r.repo.Project.UpdateByConditions(id, *data, []string{"name", "description", "sort_order"}...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProjectService) DeleteProject(userLogin string, id string) (*string, error) {
	data := r.repo.Project.FindById(id)
	if data == nil {
		return nil, fmt.Errorf("project not found")
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	data.UseYn = &STATUS_IN_ACTIVE
	data.UpdatedAt = &timeNow
	data.UpdatedBy = &userLogin
	err := r.repo.Project.UpdateByConditions(id, *data, []string{"use_yn"}...)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProjectService) ListProject(request *model.ListProjectRequest) ([]*model.Project, error) {
	result := r.repo.Project.GetList(request)
	return result, nil
}

func (r *ProjectService) GetProjectByID(id string) (*model.Project, error) {
	result := r.repo.Project.FindById(id)
	return result, nil
}
