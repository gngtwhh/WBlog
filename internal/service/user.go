package service

import (
	"errors"
	"time"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/gngtwhh/WBlog/pkg/utils"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Register(user *model.User) error {
	existUser, err := svc.repo.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if existUser != nil {
		return errors.New("username has been exist")
	}
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("internal error: hashing password failed")
	}
	user.Password = hashedPwd
	// TODO: avatar path -> Config/os.env
	if user.Avatar == "" {
		user.Avatar = "/static/default_avatar.png"
	}
	return svc.repo.Create(user)
}

func (svc *UserService) Login(username, password string) (*model.User, string, error) {
	user, err := svc.repo.GetByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("username or password is incorrect")
	}
	if !utils.CheckPassword(user.Password, password) {
		return nil, "", errors.New("username or password is incorrect")
	}
	// TODO: issuer should be load by Config/os.env
	token, err := utils.GenToken(user.ID, user.Username, user.Role, time.Hour*24, "WBLOG")
	if err != nil {
		return nil, "", err
	}
	user.Password = ""
	return user, token, nil
}

func (svc *UserService) GetProfile(id uint64) (user *model.User, err error) {
	user, err = svc.repo.GetByID(id)
	if user != nil {
		user.Password = ""
	}
	return
}

func (svc *UserService) UpdateProfile(inputUser *model.User) error {
	user, err := svc.repo.GetByID(inputUser.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not fount")
	}

	if inputUser.Nickname != "" {
		user.Nickname = inputUser.Nickname
	}
	// TODO: file upload/link check
	if inputUser.Avatar != "" {
		user.Avatar = inputUser.Avatar
	}
	return svc.repo.Update(user)
}

func (svc *UserService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	user, err := svc.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not fount")
	}

	if !utils.CheckPassword(user.Password, oldPassword) {
		return errors.New("invalid old password")
	}
	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to encrypt password")
	}
	user.Password = hashedPwd

	return svc.repo.Update(user)
}
