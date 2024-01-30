package app

import (
	"example/internal/domain"
)

type UserService interface {
	SignUp(user domain.User) error
	LoginUser(userName, pass string) (bool, error)
	FindAll(page, limit int) ([]*domain.User, error)
}

func NewUserService(repo domain.UserRespository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo domain.UserRespository
}

func (u *userService) LoginUser(userName, pass string) (bool, error) {
	ok := true

	ok, err := u.repo.GetUserName(&userName)
	if !ok || err != nil {
		return false, domain.ErrUserNotFound
	}

	ok, err = u.repo.GetPassword(&pass)
	if !ok || err != nil {
		return false, domain.ErrUserNotFound
	}

	return true, nil
}

func (u *userService) SignUp(user domain.User) error {
	err := Checker(user)
	if err != nil {
		return err
	}

	_, err = u.repo.GetUserName(&user.FullName)
	if err != nil {

		return domain.ErrUserAlreadyExists
	}

	if err := u.repo.Save(&user); err != nil {
		return err
	}
	return nil
}

func (u *userService) FindAll(page, limit int) ([]*domain.User, error) {
	return u.repo.FindAll(page, limit)
}

func Checker(user domain.User) error {
	if user.Password == "" {
		return domain.ErrInvalidPassword
	}
	if user.FullName == "" {
		return domain.ErrInvalidFullName
	}

	if user.UserName == "" {
		return domain.ErrInvalidUserName
	}
	return nil
}
