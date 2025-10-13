package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetCommentRoute(router fiber.Router, db *gorm.DB) {
	commentRoute := router.Group("/:postId/comments", middleware.AuthMiddlware())

	commentRepository := repository.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepository, db)
	commentHandler := handler.NewCommentHandler(commentService)

	commentRoute.Get("/", commentHandler.FindAllComment)
	commentRoute.Post("/", commentHandler.CreateComment)

}
