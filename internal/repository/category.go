package repository

import (
	"errors"

	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindById(id int) (*models.Category, error)
	Create(data *models.Category) error
	Update(id int, data map[string]interface{}) error
	DeleteById(id int) error
}

type CategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		DB: DB,
	}
}

func (r *CategoryRepositoryImpl) FindAll() ([]models.Category, error) {
	var categories []models.Category

	resTx := r.DB.Find(&categories)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return categories, nil
}

func (r *CategoryRepositoryImpl) FindById(id int) (*models.Category, error) {
	var user models.Category
	resTx := r.DB.Where("id = ?", id).Take(&user)

	if resTx.Error != nil {
		return nil, exception.NewGormDBErr(resTx.Error)
	}

	return &user, nil
}

func (r *CategoryRepositoryImpl) Create(data *models.Category) error {
	resTx := r.DB.Create(data)

	if resTx.Error != nil {
		return exception.NewGormDBErr(resTx.Error)
	}

	return nil
}

func (r *CategoryRepositoryImpl) Update(id int, data map[string]interface{}) error {
	resTx := r.DB.Model(&models.Category{}).Where("id = ?", id).Updates(data)

	if resTx.Error != nil {
		return exception.NewGormDBErr(resTx.Error)
	}

	if resTx.RowsAffected == 0 {
		return exception.NewGormDBErr(errors.New("no row affected"))
	}

	return nil
}

func (r *CategoryRepositoryImpl) DeleteById(id int) error {
	rsTx := r.DB.Where("id = ?", id).Delete(&models.Category{})

	if rsTx.Error != nil {
		return exception.NewGormDBErr(rsTx.Error)
	}
	return nil
}
