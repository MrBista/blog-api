package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetAuthRoute(router fiber.Router, db *gorm.DB) {

	authRepository := repository.NewUserRepository(db)
	authService := services.NewAutService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	authRoute := router.Group("/auth")

	authRoute.Post("/login", authHandler.LoginUser)
	authRoute.Post("/register", authHandler.RegisterUser)

	authRoute.Post("/google/url", authHandler.GetGoogleAuthURL)
	authRoute.Post("/google/callback", authHandler.GoogleCallback)

}
