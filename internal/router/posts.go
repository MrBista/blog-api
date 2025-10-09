package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPostRoute(router fiber.Router, db *gorm.DB) {
	postRepository := repository.NewPostRepository(db)
	postService := services.NewPostService(postRepository)
	handlerPost := handler.NewHandlerPost(postService)

	postRouter := router.Group("/posts")
	postRouter.Get("/", handlerPost.GetAllPosts)
	postRouter.Get("/:slug", handlerPost.GetPostBySlug)

}
