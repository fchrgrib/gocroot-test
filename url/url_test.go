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

func TestWebWithoutCookies(t *testing.T) {
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
