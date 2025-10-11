package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCategoryRouter(route fiber.Router, db *gorm.DB) {
	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository, db)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	categoryRouter := route.Group("/categories")

	categoryRouter.Get("/", categoryHandler.FindAllCategory)
	categoryRouter.Get("/:id", categoryHandler.FindCategoryById)
	categoryRouter.Post("/", categoryHandler.CreateCategory)

}
