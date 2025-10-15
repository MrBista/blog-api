package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPostRoute(router fiber.Router, db *gorm.DB) {
	postStroage := services.NewLocalStorage("./public", "/public")
	postRepository := repository.NewPostRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	postService := services.NewPostService(postRepository, categoryRepository, postStroage)
	handlerPost := handler.NewHandlerPost(postService)

	postRouter := router.Group("/posts")

	postRouter.Post("/uploads", handlerPost.SaveFileTemp)
	postRouter.Get("/", handlerPost.GetAllPosts)
	postRouter.Get("/:slug", middleware.AuthMiddlware(), handlerPost.GetPostBySlug)
	postRouter.Delete("/:slug", middleware.AuthMiddlware(), handlerPost.DeletePost)
	postRouter.Post("/", middleware.AuthMiddlware(), handlerPost.CreatePost)
	postRouter.Put("/:slug", middleware.AuthMiddlware(), handlerPost.UpdatePost)

	SetCommentRoute(postRouter, db)

}
