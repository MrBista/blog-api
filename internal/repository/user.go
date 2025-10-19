package repository

import (
	"fmt"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByIdentifier(identifier string) (*models.User, error)
	FindById(id int) (*models.User, error)
	FindByEmailOrUsername(email string, username string) (*models.User, error)
	FindAllUser() ([]models.User, error)
	FindAllUserWithPagination(filter dto.UserFilterRequest) (*dto.PaginationResult, error)
	DeactiveUsers(ids []int) error
	CreateFollower(follower *models.Follower) error
	DeleteFollower(followingId int, userId int) error
	GetListFollowing(userId int) ([]dto.UserFollowingDTO, error)
	GetListFollower(userId int) ([]dto.UserFollowerDTO, error)
	CountFollowing(userId int) (int64, error)
	CheckIsFollowing(followerId int, followingId int) (bool, error)
	CountFollower(userId int) (int64, error)
	GetDetailUser(userId int) (*dto.UserResponse, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	resTx := r.DB.Create(&user)

	if resTx.Error != nil {
		return fmt.Errorf("failed insert user %w", resTx.Error)
	}

	return nil
}

func (r *UserRepositoryImpl) FindByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	resTx := r.DB.Where("email = ?", identifier).Or("username = ?", identifier).First(&user)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindById(id int) (*models.User, error) {
	var user models.User
	resTx := r.DB.Where("id = ?", id).Take(&user)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmailOrUsername(email string, username string) (*models.User, error) {
	var user models.User
	resTx := r.DB.Where("email = ?", email).Or("username = ?", username).First(&user)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindAllUser() ([]models.User, error) {
	var users []models.User

	resTx := r.DB.Find(&users)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return users, nil
}

func (r *UserRepositoryImpl) FindAllUserWithPagination(filter dto.UserFilterRequest) (*dto.PaginationResult, error) {
	var users []models.User
	var total int64

	query := r.DB.Model(&models.User{})

	if filter.Email != "" {
		query.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	if filter.Username != "" {
		query.Where("username LIKE ?", "%"+filter.Username, "%")
	}

	if filter.Role != 0 {
		query.Where("role = ?", filter.Role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	query = query.Offset(filter.GetOffset()).Limit(filter.PageSize)

	if filter.Sort != "" {
		query = query.Order(filter.Sort)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return dto.NewPaginationResult(users, total, filter.Page, filter.PageSize, "users"), nil
}

func (r *UserRepositoryImpl) DeactiveUsers(ids []int) error {
	resTx := r.DB.Model(&models.User{}).
		Where("id in ?", ids).
		Update("status", 0)

	if resTx.Error != nil {
		return exception.NewGormDBErr(resTx.Error)
	}

	return nil
}

func (r *UserRepositoryImpl) CreateFollower(follower *models.Follower) error {
	if err := r.DB.Create(&follower).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	return nil
}

func (r *UserRepositoryImpl) DeleteFollower(followingId int, userId int) error {
	if err := r.
		DB.
		Where("following_id = ?", followingId).
		Where("follower_id = ?", userId).
		Delete(&models.Follower{}).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetDetailUser(userId int) (*dto.UserResponse, error) {
	var userDetail dto.UserResponse

	if err := r.
		DB.
		Table("users").
		Select("name", "username", "email", "bio", "role").
		Where("id = ?", userId).
		Scan(&userDetail).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return &userDetail, nil
}

func (r *UserRepositoryImpl) GetListFollower(userId int) ([]dto.UserFollowerDTO, error) {
	var followers []dto.UserFollowerDTO

	err := r.DB.
		Table("followers").
		Select("users.id, users.name, users.username, users.email, users.profile_image_uri, users.bio, followers.created_at as followed_at").
		Joins("JOIN users ON users.id = followers.follower_id").
		Where("followers.following_id = ?", userId).
		Scan(&followers).Error

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return followers, nil
}

func (r *UserRepositoryImpl) GetListFollowing(userId int) ([]dto.UserFollowingDTO, error) {
	var following []dto.UserFollowingDTO

	err := r.DB.
		Table("followers").
		Select("users.id, users.name, users.username, users.email, users.profile_image_uri, users.bio, followers.created_at as followed_at").
		Joins("JOIN users ON users.id = followers.following_id").
		Where("followers.follower_id = ?", userId).
		Scan(&following).Error

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return following, nil
}

func (r *UserRepositoryImpl) CountFollower(userId int) (int64, error) {
	var count int64

	err := r.DB.
		Model(&models.Follower{}).
		Where("following_id = ?", userId).
		Count(&count).Error

	if err != nil {
		return 0, exception.NewGormDBErr(err)
	}

	return count, nil
}

func (r *UserRepositoryImpl) CountFollowing(userId int) (int64, error) {
	var count int64

	err := r.DB.
		Model(&models.Follower{}).
		Where("follower_id = ?", userId).
		Count(&count).Error

	if err != nil {
		return 0, exception.NewGormDBErr(err)
	}

	return count, nil
}

func (r *UserRepositoryImpl) CheckIsFollowing(followerId int, followingId int) (bool, error) {
	var count int64

	err := r.DB.
		Model(&models.Follower{}).
		Where("follower_id = ? AND following_id = ?", followerId, followingId).
		Count(&count).Error

	if err != nil {
		return false, exception.NewGormDBErr(err)
	}

	return count > 0, nil
}
