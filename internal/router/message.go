package router

import (
	"github.com/MrBista/blog-api/internal/handler"
	"github.com/MrBista/blog-api/internal/middleware"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetMessageRouter(router fiber.Router, db *gorm.DB) {

	chatRoomRepository := repository.NewChatRoom(db)
	chatMessageRepository := repository.NewChatMessageRepository(db)

	messageService := services.NewMessageServiceImpl(chatRoomRepository, chatMessageRepository)

	messageHandler := handler.NewMessageServiceImpl(messageService)

	chatRoomRouter := router.Group("api/v1/chat/room")
	chatRoomRouter.Get("/", middleware.AuthMiddlware(), messageHandler.FindOrCreateRoom)
}
