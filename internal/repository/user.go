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
