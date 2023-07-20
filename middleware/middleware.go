package middleware

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/pisato"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Middleware(c *fiber.Ctx) error {
	token := c.Cookies("token")

	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status": "Unauthorized",
		})
	}

	maker, err := pisato.NewPasetoMaker(config.PrivateKey)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	_, err = maker.VerifyToken(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(err)
	}

	return c.Next()
}
