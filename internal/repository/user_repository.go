package repository

import (
	"tsv-golang/internal/dto"
	"tsv-golang/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	GetList(param *dto.GetListUserRequest) []*entity.User
}

func UserRepositoryInit(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ UserRepositoryInterface = &UserRepository{}

func (repo UserRepository) GetList(param *dto.GetListUserRequest) []*entity.User {
	result := make([]*entity.User, 0)
	query := repo.db.Table("tb_user").Select("*")
	if param.Id != "" {
		query = query.Where("id = ?", param.Id)
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.Limit == 0 {
		param.Limit = 20
	}
	query.Offset(int(param.Page-1) * int(param.Limit)).Limit(int(param.Limit)).Order("updated_at desc")
	query = query.Find(&result)
	return result
}
