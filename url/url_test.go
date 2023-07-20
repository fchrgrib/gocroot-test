package url

import (
	"github.com/gocroot/gocroot/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {
	app := fiber.New()

	app.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	app.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

	app.Get("/", controller.Homepage) //ujicoba panggil package musik

	req := httptest.NewRequest("POST", "/api/whatsauth/request", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req = httptest.NewRequest("GET", "/", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req = httptest.NewRequest("GET", "/ws/whatsauth/qr", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 426, resp.StatusCode)
}

func TestAuthSession(t *testing.T) {
	app := fiber.New()

	app.Post("/login", controller.LoginUser)
	app.Post("/register", controller.RegisterUser)
	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("token")
		return nil
	})

	req := httptest.NewRequest("POST", "/login", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req = httptest.NewRequest("POST", "/register", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req = httptest.NewRequest("POST", "/logout", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 405, resp.StatusCode)
}
