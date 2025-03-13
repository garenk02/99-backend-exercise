// usecase/user_usecase.go
package usecase

import (
	"public-api/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) GetUserByID(id int) (*domain.User, error) {
	return u.userRepo.GetUserByID(id)
}

func (u *UserUseCase) GetUsers(pageNum, pageSize int) ([]*domain.User, error) {
	return u.userRepo.GetUsers(pageNum, pageSize)
}

func (u *UserUseCase) CreateUser(name string) (*domain.User, error) {
	return u.userRepo.CreateUser(name)
}
