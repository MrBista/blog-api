package router

import (
	"github.com/MrBista/blog-api/internal/database"
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupAllRoutes(app *fiber.App) {
	router := app.Group("/api")

	userRepository := repository.NewUserRepository(database.DB)
	subscriptionService := services.NewXenditPaymentService(userRepository, database.DB)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	router.Post("/webhook/xendit", subscriptionHandler.WebhookPayment)
	subscription := router.Group("/subscriptions", middleware.AuthMiddlware())

	subscription.Post("/", subscriptionHandler.CreateSubscription)

	SetupPostRoute(router, database.DB)
	SetAuthRoute(router, database.DB)
	SetupCategoryRouter(router, database.DB)
	SetUserRoute(router, database.DB)
	// SetCommentRoute(router, database.DB)
}
