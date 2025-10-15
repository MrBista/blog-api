package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUploadFile(router fiber.Router, db *gorm.DB) {

	router.Post("/upload", func(c *fiber.Ctx) error {
		// di multipartnya ambil apakah ini type post
		// kalau dia type post maka simpan

		return nil

	})

}
