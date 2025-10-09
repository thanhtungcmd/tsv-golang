package service

import (
	"tsv-golang/internal/graph/model"
	"tsv-golang/internal/repository"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.Repositories
}

type UserServiceInterface interface {
	CreateUser(input model.UserInput) (*model.User, error)
}

func UserServiceInit(repo *repository.Repositories) *UserService {
	return &UserService{
		repo: repo,
	}
}

var _ UserServiceInterface = &UserService{}

func (u *UserService) CreateUser(input model.UserInput) (*model.User, error) {
	user := &model.User{}
	err := mapstructure.Decode(input, &user)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	data, err := u.repo.User.CreateAndReturn(user)
	if err != nil {
		return nil, err
	}
	return data, nil
}
