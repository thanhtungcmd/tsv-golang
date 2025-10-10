package service

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"
	"tsv-golang/internal/graph/model"
	"tsv-golang/internal/repository"
	"tsv-golang/pkg/datetime"
	"tsv-golang/pkg/mail"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.Repositories
}

var userActive = 1
var userInActive = 0

type UserServiceInterface interface {
	CreateUser(userLogin string, input model.UserInput) (*model.User, error)
	ListUsers(request *model.ListUsersRequest) ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUser(userLogin string, id string, input model.UserUpdateInput) (*model.User, error)
	Login(username string, password string) (*model.LoginResponse, error)
	ForgetPassword(email string) (*string, error)
}

func UserServiceInit(repo *repository.Repositories) *UserService {
	return &UserService{
		repo: repo,
	}
}

var _ UserServiceInterface = &UserService{}

func (u *UserService) Login(username string, password string) (*model.LoginResponse, error) {
	result := &model.LoginResponse{}
	user := u.repo.User.FindByUsername(username)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "defaultSecretKey"
	}
	claims := jwt.MapClaims{
		"sub":      username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
		"username": username,
		"email":    user.Email,
		"phone":    user.PhoneNumber,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["typ"] = "JWT"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	result.AccessToken = tokenString
	result.User = user
	return result, nil
}

func (u *UserService) CreateUser(userLogin string, input model.UserInput) (*model.User, error) {
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
	timeNow := datetime.Datetime().TimeNow().ToString()
	user.CreatedAt = &timeNow
	user.UpdatedAt = &timeNow
	user.CreatedBy = &userLogin
	user.UpdatedBy = &userLogin
	user.UseYn = &userActive
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

func (u *UserService) UpdateUser(userLogin string, id string, input model.UserUpdateInput) (*model.User, error) {
	user := u.repo.User.FindById(id)
	err := mapstructure.Decode(input, &user)
	if err != nil {
		return nil, err
	}
	timeNow := datetime.Datetime().TimeNow().ToString()
	user.UpdatedAt = &timeNow
	user.UpdatedBy = &userLogin
	err = u.repo.User.UpdateById(id, *user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) ForgetPassword(email string) (*string, error) {
	user := u.repo.User.FindByEmail(email)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	// Save Code
	code := fmt.Sprintf("%06d", rand.IntN(1_000_000))
	user.VerifyCode = &code
	err := u.repo.User.UpdateByConditions(user.ID, *user, []string{"verify_code"}...)
	if err != nil {
		return nil, err
	}
	// Send mail
	to := make([]string, 0)
	to = append(to, email)
	err = mail.SendEmail(to, "Quên mật khẩu", fmt.Sprintf("Verify Code: %s", code))
	if err != nil {
		return nil, err
	}
	return &email, nil
}
