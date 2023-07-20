package url

import (
	"github.com/gocroot/gocroot/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebWithoutCookies(t *testing.T) {
	t.Helper()
	app := fiber.New()

	app.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	app.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

	app.Get("/", controller.Homepage) //ujicoba panggil package musik
	app.Get("/presensi", controller.GetPresensiBulanIni)

	route := []string{"/", "/presensi"}
	request := []string{"GET", "GET"}

	for i := 0; i < len(route); i++ {
		req := httptest.NewRequest(request[i], route[i], nil)

		resp, _ := app.Test(req, 3)
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "expected 200")
	}
}
