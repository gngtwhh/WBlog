package service

import (
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Create(user *model.User) error {
	return nil
}
func (svc *UserService) GetByUsername(username string) (user *model.User, err error) {
	return nil, nil
}
func (svc *UserService) GetByID(id uint64) (user *model.User, err error) {
	return nil, nil
}
func (svc *UserService) Update(user *model.User) error {
	return nil
}
