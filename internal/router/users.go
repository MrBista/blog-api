package router

import (
	"github.com/MrBista/blog-api/internal/enum"
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUserRoute(router fiber.Router, db *gorm.DB) {
	userRoute := router.Group("/users")

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository, db)
	userHandler := handler.NewUserHandler(userService)

	userRoute.Get("/", middleware.AuthMiddlware(), userHandler.GetAllUser)
	userRoute.Post("/", middleware.AuthMiddlware(), middleware.RoleMiddleare(enum.RoleAdmin), userHandler.CreateUser)
	userRoute.Put("/deactive", middleware.AuthMiddlware(), middleware.RoleMiddleare(enum.RoleAdmin), userHandler.DeactiveUser)
}
