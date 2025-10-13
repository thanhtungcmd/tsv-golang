package repository

import (
	"tsv-golang/internal/graph/model"
	"tsv-golang/pkg/ulid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectRepository struct {
	db *gorm.DB
}

type ProjectRepositoryInterface interface {
	CreateAndReturn(model *model.Project) (*model.Project, error)
	GetList(param *model.ListProjectRequest) []*model.Project
	FindById(id string) *model.Project
	UpdateByConditions(id string, model model.Project, fieldsUpdate ...string) error
}

func ProjectRepositoryInit(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

var _ ProjectRepositoryInterface = &ProjectRepository{}

func (repo ProjectRepository) CreateAndReturn(model *model.Project) (*model.Project, error) {
	model.ID = ulid.NewULID()
	result := repo.db.Table("balheh.tb_projects").Clauses(clause.Returning{}).Create(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo ProjectRepository) GetList(param *model.ListProjectRequest) []*model.Project {
	result := make([]*model.Project, 0)
	query := repo.db.Table("balheh.tb_projects").Select("*")
	if param != nil {
		if param.Offset != nil {
			param.Offset = &defaultOffset
		}
		if param.Limit != nil {
			param.Limit = &defaultLimit
		}
		if param.Offset != nil && param.Limit != nil {
			query.Offset(*param.Offset).Limit(*param.Limit)
		}
	}
	query.Order("updated_at, sort_order desc")
	query = query.Find(&result)
	return result
}

func (repo ProjectRepository) FindById(id string) *model.Project {
	var user *model.Project
	item := repo.db.Table("balheh.tb_projects").Take(&user, "id = ?", id)
	if item.RowsAffected == 0 {
		return nil
	}
	return user
}

func (repo ProjectRepository) UpdateByConditions(id string, model model.Project, fieldsUpdate ...string) error {
	item := repo.db.Table("balheh.tb_projects").Where("id = ?", id).Select(fieldsUpdate).Updates(model)
	if item.RowsAffected == 0 {
		return item.Error
	}
	return nil
}
