package middleware

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/pisato"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareNoToken(t *testing.T) {

	app := fiber.New()

	app.Use(Middleware)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})
	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestMiddlewareValidToken(t *testing.T) {
	app := fiber.New()

	maker, _ := pisato.NewPasetoMaker(config.PrivateKey)
	token, _ := maker.CreateToken("fahrian", 4*time.Minute)

	app.Use(Middleware)

	app.Get("/this", func(ctx *fiber.Ctx) error {
		require.NotEmpty(t, ctx.Cookies("token"))
		return ctx.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/this", nil)
	req.Header.Set("Cookie", "token="+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewareInvalidToken(t *testing.T) {
	app := fiber.New()

	app.Use(Middleware)

	app.Get("/this", func(ctx *fiber.Ctx) error {
		require.NotEmpty(t, ctx.Cookies("token"))
		return ctx.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/this", nil)
	req.Header.Set("Cookie", "token=invalid token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
