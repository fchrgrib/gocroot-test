package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorLogRes(statusCode int, err error, statusDesc string, c *fiber.Ctx) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status": statusDesc,
		"error":  err,
	})
}
