package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/gngtwhh/WBlog/internal/cache"
	"github.com/gngtwhh/WBlog/internal/config"
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/gngtwhh/WBlog/pkg/utils"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("username already exists")
	ErrAuthFailed     = errors.New("username or password is incorrect")
	ErrInvalidOldPass = errors.New("invalid old password")
)

type UserService struct {
	repo repository.UserRepository
	log  *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  logger.With("component", "user_service"),
	}
}

func (svc *UserService) Register(user *model.User) error {
	existUser, err := svc.repo.GetByUsername(user.Username)
	if err == nil {
		if existUser != nil {
			return ErrUserExists
		}
		svc.log.Error("failed to check username existence", "err", err)
		return err
	}
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		svc.log.Error("failed to hash password", "err", err)
		return errors.New("internal error: hashing password failed")
	}
	user.Password = hashedPwd
	// TODO: avatar path -> Config/os.env
	if user.Avatar == "" {
		user.Avatar = "/static/default_avatar.png"
	}
	if err := svc.repo.Create(user); err != nil {
		svc.log.Error("failed to create user", "username", user.Username, "err", err)
		return err
	}
	return nil
}

func (svc *UserService) Login(username, password string) (*model.User, string, error) {
	user, err := svc.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrAuthFailed
		}
		svc.log.Error("login failed: db query error", "err", err)
		return nil, "", err
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, "", ErrAuthFailed
	}
	// TODO: issuer should be load by Config/os.env
	token, err := utils.GenToken(user.ID, user.Username, user.Role, config.Cfg.GetJwtDuration(), "WBLOG")
	if err != nil {
		svc.log.Error("failed to generate token", "uid", user.ID, "err", err)
		return nil, "", err
	}
	user.Password = ""
	return user, token, nil
}

func (svc *UserService) GetProfile(id uint64) (*model.User, error) {
	user, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		svc.log.Error("failed to get user profile", "uid", id, "err", err)
		return nil, err
	}
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (svc *UserService) UpdateProfile(inputUser *model.User) error {
	user, err := svc.repo.GetByID(inputUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	needUpdate := false
	if inputUser.Nickname != "" && inputUser.Nickname != user.Nickname {
		user.Nickname = inputUser.Nickname
		needUpdate = true
	}
	// TODO: file upload/link check
	if inputUser.Avatar != "" && inputUser.Avatar != user.Avatar {
		user.Avatar = inputUser.Avatar
		needUpdate = true
	}
	if !needUpdate {
		return nil
	}
	if err := svc.repo.Update(user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		svc.log.Error("failed to update profile", "uid", user.ID, "err", err)
		return err
	}
	return nil
}

func (svc *UserService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	user, err := svc.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if !utils.CheckPassword(user.Password, oldPassword) {
		return ErrInvalidOldPass
	}
	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		svc.log.Error("faied to hash new password", "uid", userID, "err", err)
		return errors.New("internal error: failed to hash password")
	}
	user.Password = hashedPwd

	if err := svc.repo.Update(user); err != nil {
		svc.log.Error("failed to update password", "uid", user.ID, "err", err)
		return err
	}
	return nil
}

func (svc *UserService) Logout(tokenStr string, exp int64) error {
	now := time.Now()
	expTime := time.Unix(exp, 0)
	if now.After(expTime) {
		return nil
	}
	duration := expTime.Sub(now)

	key := cache.PrefixJWTBlacklist + tokenStr
	err := cache.RDB.Set(context.Background(), key, "1", duration).Err()
	if err != nil {
		svc.log.Error("failed to add token to blacklist", "key", key, "err", err)
		return err
	}
	return nil
}
