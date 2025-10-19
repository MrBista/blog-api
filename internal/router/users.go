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

	// My followers & following (harus di atas /:id agar tidak bentrok)
	userRoute.Get("/me/followers", middleware.AuthMiddlware(), userHandler.GetMyFollowers)
	userRoute.Get("/me/following", middleware.AuthMiddlware(), userHandler.GetMyFollowing)

	userRoute.Get("/:id", middleware.AuthMiddlware(), userHandler.GetDetailUser)

	// Follow/Unfollow user
	userRoute.Post("/:id/follow", middleware.AuthMiddlware(), userHandler.FollowUser)
	userRoute.Delete("/:id/follow", middleware.AuthMiddlware(), userHandler.UnfollowUser)

	// Check follow status
	userRoute.Get("/:id/follow/status", middleware.AuthMiddlware(), userHandler.CheckFollowStatus)

	// Get followers & following of specific user
	userRoute.Get("/:id/followers", middleware.AuthMiddlware(), userHandler.GetListFollower)
	userRoute.Get("/:id/following", middleware.AuthMiddlware(), userHandler.GetListFollowing)

	// Count followers & following
	userRoute.Get("/:id/followers/count", middleware.AuthMiddlware(), userHandler.GetFollowerCount)
	userRoute.Get("/:id/following/count", middleware.AuthMiddlware(), userHandler.GetFollowingCount)
}
