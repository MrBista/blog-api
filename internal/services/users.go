package services

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	FindAllUsers() ([]models.User, error)
	FindAllUserWithPaginatin(filter dto.UserFilterRequest) (*dto.PaginationResult, error)
	DeactiveUsers(ids []int) error
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *gorm.DB
}

func NewUserService(userRepository repository.UserRepository, db *gorm.DB) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (s *UserServiceImpl) FindAllUsers() ([]models.User, error) {
	datas, err := s.UserRepository.FindAllUser()

	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (s *UserServiceImpl) FindAllUserWithPaginatin(filter dto.UserFilterRequest) (*dto.PaginationResult, error) {
	datas, err := s.UserRepository.FindAllUserWithPagination(filter)

	if err != nil {
		return nil, err
	}

	return datas, err
}

func (s *UserServiceImpl) DeactiveUsers(ids []int) error {
	if err := s.UserRepository.DeactiveUsers(ids); err != nil {
		return err
	}

	return nil
}
