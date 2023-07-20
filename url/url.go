package url

import (
	"github.com/gocroot/gocroot/controller"
	"github.com/gocroot/gocroot/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Web(page *fiber.App) {
	page.Use(middleware.Middleware)
	page.Post("/api/whatsauth/request", controller.PostWhatsAuthRequest)  //API from user whatsapp message from iteung gowa
	page.Get("/ws/whatsauth/qr", websocket.New(controller.WsWhatsAuthQR)) //websocket whatsauth

	page.Get("/", controller.Homepage) //ujicoba panggil package musik
	page.Get("/presensi", controller.GetPresensiBulanIni)
}

func AuthSession(page *fiber.App) {
	page.Post("/login", controller.LoginUser)
	page.Post("/register", controller.RegisterUser)
	page.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("token")
		return nil
	})
}
