package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupReadingListRoutes(app fiber.Router, postService services.PostService) {
	readingListHandler := handler.NewReadingListHandler(postService) // !TODO pindah post service ke service tersendiri

	// Group route dengan prefix /api/reading-lists
	readingList := app.Group("/reading-lists", middleware.AuthMiddlware())

	// Reading List Management
	readingList.Post("/", readingListHandler.CreateReadingList)
	readingList.Get("/", readingListHandler.GetReadingLists)
	readingList.Get("/:id", readingListHandler.GetReadingListByID)
	readingList.Put("/:id", readingListHandler.UpdateReadingList)
	readingList.Delete("/:id", readingListHandler.DeleteReadingList)

	// Saved Posts Management
	readingList.Post("/saved-posts", readingListHandler.CreateSavedPost)
	readingList.Get("/:listId/saved-posts", readingListHandler.GetSavedPosts)
	readingList.Put("/saved-posts/:id", readingListHandler.UpdateSavedPost)
	readingList.Delete("/saved-posts/:id", readingListHandler.DeleteSavedPost)
	readingList.Delete("/saved-posts", readingListHandler.DeleteSavedPostByPostAndList)

	// Mark all as read
	readingList.Post("/:listId/mark-all-read", readingListHandler.MarkAllAsRead)
}
