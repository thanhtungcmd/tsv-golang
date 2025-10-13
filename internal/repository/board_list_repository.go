package repository

import (
	"tsv-golang/internal/graph/model"
	"tsv-golang/pkg/ulid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BoardListRepository struct {
	db *gorm.DB
}

type BoardListRepositoryInterface interface {
	CreateAndReturn(model *model.BoardList) (*model.BoardList, error)
	GetList(param *model.ListBoardListRequest) []*model.BoardList
	FindById(id string) *model.BoardList
	UpdateByConditions(id string, model model.BoardList, fieldsUpdate ...string) error
}

func BoardListRepositoryInit(db *gorm.DB) *BoardListRepository {
	return &BoardListRepository{
		db: db,
	}
}

var _ BoardListRepositoryInterface = &BoardListRepository{}

func (repo BoardListRepository) CreateAndReturn(model *model.BoardList) (*model.BoardList, error) {
	model.ID = ulid.NewULID()
	result := repo.db.Table("balheh.tb_board_lists").Clauses(clause.Returning{}).Create(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo BoardListRepository) GetList(param *model.ListBoardListRequest) []*model.BoardList {
	result := make([]*model.BoardList, 0)
	query := repo.db.Table("balheh.tb_board_lists").Select("*")
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

func (repo BoardListRepository) FindById(id string) *model.BoardList {
	var user *model.BoardList
	item := repo.db.Table("balheh.tb_board_lists").Take(&user, "id = ?", id)
	if item.RowsAffected == 0 {
		return nil
	}
	return user
}

func (repo BoardListRepository) UpdateByConditions(id string, model model.BoardList, fieldsUpdate ...string) error {
	item := repo.db.Table("balheh.tb_board_lists").Where("id = ?", id).Select(fieldsUpdate).Updates(model)
	if item.RowsAffected == 0 {
		return item.Error
	}
	return nil
}
