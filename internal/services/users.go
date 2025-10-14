package services

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/enum"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	FindAllUsers() ([]models.User, error)
	FindAllUserWithPaginatin(filter dto.UserFilterRequest) (*dto.PaginationResult, error)
	DeactiveUsers(ids []int) error
	CreateUser(userBody dto.RegisterRequest) (*dto.UserResponse, error)
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

func (s *UserServiceImpl) FindByEmailOrUsername(email, username string) (*models.User, error) {
	user, err := s.UserRepository.FindByEmailOrUsername(email, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) CreateUser(userBody dto.RegisterRequest) (*dto.UserResponse, error) {

	_, err := s.FindByEmailOrUsername(userBody.Email, userBody.Username)

	if err == nil {
		return nil, exception.NewBadRequestErr("Username/Password is invalid")
	}

	passwordHash, err := utils.HashPassword(userBody.Password)
	if err != nil {
		return nil, err
	}

	if !enum.IsValidRole(enum.UserRole(userBody.Role)) {
		return nil, exception.NewBadRequestErr("Invalid role value")
	}

	modelUser := models.User{
		Name:     userBody.Email,
		Username: userBody.Username,
		Email:    userBody.Email,
		Password: passwordHash,
		Status:   1, // 1 == aktif,
		Role:     userBody.Role,
	}

	if err := s.UserRepository.CreateUser(&modelUser); err != nil {
		return nil, err
	}

	mapUserResponse := dto.UserResponse{
		Id:       int(modelUser.ID),
		Name:     modelUser.Name,
		Username: modelUser.Username,
		Email:    modelUser.Email,
		Role:     modelUser.Role,
	}

	if modelUser.Bio != nil {
		mapUserResponse.Bio = *modelUser.Bio
	}
	return &mapUserResponse, nil
}
