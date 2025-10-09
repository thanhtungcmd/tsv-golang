package repository

import (
	"errors"
	"math/rand"
	"time"
	"tsv-golang/internal/graph/model"

	"github.com/oklog/ulid/v2"
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
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	model.ID = id.String()
	result := repo.db.Table("balheh.tb_user").Clauses(clause.Returning{}).Create(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo UserRepository) Login(username string) (*model.User, error) {
	var user model.User
	if err := repo.db.Table("balheh.tb_user").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (repo UserRepository) GetList(param *model.ListUsersRequest) []*model.User {
	result := make([]*model.User, 0)
	query := repo.db.Table("balheh.tb_user").Select("*")
	if param.ID != nil {
		query = query.Where("id = ?", param.ID)
	}
	if param.Offset != nil {
		param.Offset = &defaultOffset
	}
	if param.Limit != nil {
		param.Limit = &defaultLimit
	}
	query.Offset(*param.Offset).Limit(*param.Limit).Order("updated_at desc")
	query = query.Find(&result)
	return result
}
