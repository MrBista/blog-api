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
	FollowUser(userToFollow int, userDetail *utils.Claims) error
	UnFollowUser(userToUnFollow int, userDetail *utils.Claims) error
	GetListFollower(userId int) ([]dto.UserFollowerDTO, error)
	GetListFollowing(userId int) ([]dto.UserFollowingDTO, error)
	CountFollower(userId int) (int64, error)
	CountFollowing(userId int) (int64, error)
	CheckIsFollowing(targetUserId int, currentUserId int) (bool, error)
	DetailUser(userId int) (*dto.UserResponse, error)
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

func (s *UserServiceImpl) DetailUser(userId int) (*dto.UserResponse, error) {
	userResponse, err := s.UserRepository.GetDetailUser(userId)

	if err != nil {
		return nil, err
	}
	return userResponse, nil
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

func (s *UserServiceImpl) FollowUser(userToFollow int, userDetail *utils.Claims) error {
	if userToFollow == userDetail.UserId {
		return exception.NewBadRequestErr("Cannot follow yourself")
	}

	var userExists int64
	if err := s.DB.Model(&models.User{}).Where("id = ?", userToFollow).Count(&userExists).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	if userExists == 0 {
		return exception.NewNotFoundErr("User not found")
	}

	isFollowing, err := s.UserRepository.CheckIsFollowing(userDetail.UserId, userToFollow)
	if err != nil {
		return err
	}
	if isFollowing {
		return exception.NewBadRequestErr("Already following this user")
	}

	followerCreate := models.Follower{
		FollowerID:  uint64(userDetail.UserId),
		FollowingID: uint64(userToFollow),
	}

	if err := s.UserRepository.CreateFollower(&followerCreate); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) UnFollowUser(userToUnFollow int, userDetail *utils.Claims) error {
	if userToUnFollow == userDetail.UserId {
		return exception.NewBadRequestErr("Cannot unfollow yourself")
	}

	isFollowing, err := s.UserRepository.CheckIsFollowing(userDetail.UserId, userToUnFollow)
	if err != nil {
		return err
	}
	if !isFollowing {
		return exception.NewNotFoundErr("Not following this user")
	}

	if err := s.UserRepository.DeleteFollower(userToUnFollow, userDetail.UserId); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) GetListFollower(userId int) ([]dto.UserFollowerDTO, error) {
	var userExists int64
	if err := s.DB.Model(&models.User{}).Where("id = ?", userId).Count(&userExists).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}
	if userExists == 0 {
		return nil, exception.NewNotFoundErr("User not found")
	}

	followers, err := s.UserRepository.GetListFollower(userId)
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (s *UserServiceImpl) GetListFollowing(userId int) ([]dto.UserFollowingDTO, error) {
	var userExists int64
	if err := s.DB.Model(&models.User{}).Where("id = ?", userId).Count(&userExists).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}
	if userExists == 0 {
		return nil, exception.NewNotFoundErr("User not found")
	}

	following, err := s.UserRepository.GetListFollowing(userId)
	if err != nil {
		return nil, err
	}

	return following, nil
}

func (s *UserServiceImpl) CountFollower(userId int) (int64, error) {
	count, err := s.UserRepository.CountFollower(userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserServiceImpl) CountFollowing(userId int) (int64, error) {
	count, err := s.UserRepository.CountFollowing(userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserServiceImpl) CheckIsFollowing(targetUserId int, currentUserId int) (bool, error) {
	if targetUserId == currentUserId {
		return false, nil
	}

	isFollowing, err := s.UserRepository.CheckIsFollowing(currentUserId, targetUserId)
	if err != nil {
		return false, err
	}

	return isFollowing, nil
}
