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
	ListUsers(request *model.ListUsersRequest) ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUser(id string, input model.UserInput) (*model.User, error)
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

func (u *UserService) ListUsers(request *model.ListUsersRequest) ([]*model.User, error) {
	result := u.repo.User.GetList(request)
	return result, nil
}

func (u *UserService) GetUserByID(id string) (*model.User, error) {
	result := u.repo.User.FindById(id)
	return result, nil
}

func (u *UserService) UpdateUser(id string, input model.UserInput) (*model.User, error) {
	user := u.repo.User.FindById(id)
	err := mapstructure.Decode(input, &user)
	if err != nil {
		return nil, err
	}
	err = u.repo.User.UpdateById(id, *user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
