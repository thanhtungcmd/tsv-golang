package repository

import (
	"tsv-golang/internal/graph/model"
	"tsv-golang/pkg/ulid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BoardRepository struct {
	db *gorm.DB
}

type BoardRepositoryInterface interface {
	CreateAndReturn(model *model.Board) (*model.Board, error)
	GetList(param *model.ListBoardRequest) []*model.Board
	FindById(id string) *model.Board
	UpdateByConditions(id string, model model.Board, fieldsUpdate ...string) error
}

func BoardRepositoryInit(db *gorm.DB) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

var _ BoardRepositoryInterface = &BoardRepository{}

func (repo BoardRepository) CreateAndReturn(model *model.Board) (*model.Board, error) {
	model.ID = ulid.NewULID()
	result := repo.db.Table("balheh.tb_boards").Clauses(clause.Returning{}).Create(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo BoardRepository) GetList(param *model.ListBoardRequest) []*model.Board {
	result := make([]*model.Board, 0)
	query := repo.db.Table("balheh.tb_boards").Select("*")
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
	query.Order("updated_at desc")
	query = query.Find(&result)
	return result
}

func (repo BoardRepository) FindById(id string) *model.Board {
	var user *model.Board
	item := repo.db.Table("balheh.tb_boards").Take(&user, "id = ?", id)
	if item.RowsAffected == 0 {
		return nil
	}
	return user
}

func (repo BoardRepository) UpdateByConditions(id string, model model.Board, fieldsUpdate ...string) error {
	item := repo.db.Table("balheh.tb_boards").Where("id = ?", id).Select(fieldsUpdate).Updates(model)
	if item.RowsAffected == 0 {
		return item.Error
	}
	return nil
}
