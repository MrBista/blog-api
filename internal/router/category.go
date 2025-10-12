package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
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

	categoryRouter.Get("/:id", middleware.AuthMiddlware(), categoryHandler.FindCategoryById)
	categoryRouter.Put("/:id", middleware.AuthMiddlware(), categoryHandler.UpdateCategory)
	categoryRouter.Delete("/:id", middleware.AuthMiddlware(), categoryHandler.DeleteCategory)
	categoryRouter.Get("/", categoryHandler.FindAllCategory)
	categoryRouter.Post("/", middleware.AuthMiddlware(), categoryHandler.CreateCategory)

}
