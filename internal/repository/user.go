package repository

import (
	"fmt"

	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByIdentifier(identifier string) (*models.User, error)
	FindById(id int) (*models.User, error)
	FindByEmailOrUsername(email string, username string) (*models.User, error)
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
