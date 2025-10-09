package repository

import (
	"errors"
	"tsv-golang/internal/graph/model"
	"tsv-golang/pkg/ulid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	GetList(param *model.ListUsersRequest) []*model.User
	CreateAndReturn(model *model.User) (*model.User, error)
	Login(username string) (*model.User, error)
	FindById(id string) *model.User
}

func UserRepositoryInit(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ UserRepositoryInterface = &UserRepository{}

var defaultOffset = 1
var defaultLimit = 20

func (repo UserRepository) CreateAndReturn(model *model.User) (*model.User, error) {
	model.ID = ulid.NewULID()
	result := repo.db.Table("balheh.tb_users").Clauses(clause.Returning{}).Create(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo UserRepository) Login(username string) (*model.User, error) {
	var user model.User
	if err := repo.db.Table("balheh.tb_users").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo UserRepository) GetList(param *model.ListUsersRequest) []*model.User {
	result := make([]*model.User, 0)
	query := repo.db.Table("balheh.tb_users").Select("*")
	if param != nil {
		if param.ID != nil {
			query = query.Where("id = ?", param.ID)
		}
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

func (repo UserRepository) FindById(id string) *model.User {
	var user *model.User
	item := repo.db.Table("balheh.tb_users").Take(&user, "id = ?", id)
	if item.RowsAffected == 0 {
		return nil
	}
	return user
}
