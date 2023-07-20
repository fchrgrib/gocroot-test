package main

import (
	"log"

	"github.com/gocroot/gocroot/config"

	"github.com/aiteung/musik"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/whatsauth/whatsauth"

	"github.com/gocroot/gocroot/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	go whatsauth.RunHub()
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.AuthSession(site)
	url.Web(site)
	log.Fatal(site.Listen(musik.Dangdut()))
}
