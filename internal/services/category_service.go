package services

import (
	"time"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"gorm.io/gorm"
)

type CategoryService interface {
	FindAllCategory() ([]dto.CategoryResponse, error)
	FindById(id int) (*dto.CategoryResponse, error)
	CreateCategory(body dto.CategoryRequst) error
	UpdateCategory(id int, body dto.CategoryRequst) error
	DeleteById(id int) error
}

type CategoryServiceImpl struct {
	DB                 *gorm.DB
	CategoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository, DB *gorm.DB) CategoryService {
	return &CategoryServiceImpl{
		DB:                 DB,
		CategoryRepository: categoryRepository,
	}
}

func (s *CategoryServiceImpl) FindAllCategory() ([]dto.CategoryResponse, error) {
	var categories []dto.CategoryResponse

	findAllCategories, err := s.CategoryRepository.FindAll()

	if err != nil {
		return nil, err
	}

	for _, category := range findAllCategories {
		categoryDto := dto.CategoryResponse{
			Id:   int(category.ID),
			Name: category.Name,
			Slug: category.Slug,
			Desc: category.Slug,
		}

		if category.ParentID != nil {
			categoryDto.ParentId = int(*category.ParentID)
		}

		categories = append(categories, categoryDto)
	}

	return categories, nil
}

func (s *CategoryServiceImpl) FindById(id int) (*dto.CategoryResponse, error) {
	category, err := s.CategoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	categoryDto := dto.CategoryResponse{
		Id:   int(category.ID),
		Name: category.Name,
		Slug: category.Slug,
		Desc: category.Slug,
	}

	if category.ParentID != nil {
		categoryDto.ParentId = int(*category.ParentID)
	}

	return &categoryDto, nil

}

func (s *CategoryServiceImpl) CreateCategory(body dto.CategoryRequst) error {
	category := models.Category{
		Name:        body.Name,
		Description: &body.Desc,
		Slug:        body.Name + time.Now().String(),
	}
	if body.ParentId != 0 {
		var parentId int64 = int64(body.ParentId)
		category.ParentID = &parentId
	}
	err := s.CategoryRepository.Create(&category)

	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryServiceImpl) UpdateCategory(id int, body dto.CategoryRequst) error {
	_, err := s.FindById(id)

	if err != nil {
		return exception.NewNotFoundErr("category not found")
	}

	categoryUpdate := map[string]interface{}{
		"Name":        body.Name,
		"Description": body.Desc,
		"ParentId":    body.ParentId,
		"Slug":        body.Name + time.Now().String(),
	}

	err = s.CategoryRepository.Update(id, categoryUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryServiceImpl) DeleteById(id int) error {
	_, err := s.FindById(id)

	if err != nil {
		return exception.NewNotFoundErr("category not found")
	}

	err = s.CategoryRepository.DeleteById(id)

	if err != nil {
		return err
	}
	return nil
}
