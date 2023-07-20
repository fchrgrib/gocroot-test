package middleware

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/pisato"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func createTestContext(token string) *fiber.Ctx {
	app := fiber.New()
	req := &fasthttp.RequestCtx{}
	req.Request.Header.SetCookie("token", token)
	return app.AcquireCtx(req)
}

func TestMiddleware_ValidToken(t *testing.T) {
	// Replace "validToken" with an actual valid token
	maker, _ := pisato.NewPasetoMaker(config.PrivateKey)
	require.NotEmpty(t, maker)

	token, _ := maker.CreateToken("fahri", 4*time.Minute)
	require.NotEmpty(t, token)

	ctx := createTestContext(token)
	require.NotNil(t, ctx)
	err := Middleware(ctx)
	require.NoError(t, err)

	//assert.Nil(t, err, "Expected error to be nil")
	//assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode(), "Expected status code 200")
}
