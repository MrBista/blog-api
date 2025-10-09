package repository

import (
	"fmt"

	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindById(id int) (*models.Category, error)
	Create(data *models.Category) error
	Update(data models.Category) error
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
	panic("not implemented") // TODO: Implement
}

func (r *CategoryRepositoryImpl) FindById(id int) (*models.Category, error) {
	var user models.Category
	resTx := r.DB.Where("id = ?", id).Take(&user)

	if resTx.Error != nil {
		return nil, fmt.Errorf("something went wrong %w", resTx.Error)
	}

	return &user, nil
}

func (r *CategoryRepositoryImpl) Create(data *models.Category) error {
	resTx := r.DB.Create(data)

	if resTx.Error != nil {
		return fmt.Errorf("something when wrong %w", resTx.Error)
	}

	return nil
}

func (r *CategoryRepositoryImpl) Update(data models.Category) error {
	panic("not implemented") // TODO: Implement
}

func (r *CategoryRepositoryImpl) DeleteById(id int) error {
	panic("not implemented") // TODO: Implement
}
