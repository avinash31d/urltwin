package main

import (
	"log"
	"os"

	"github.com/avinash31d/urltwin/config"
	"github.com/avinash31d/urltwin/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadConfig()

	app := fiber.New()

	if os.Getenv("ENV") != "production" {
		app.Use(func(c *fiber.Ctx) error {
			c.Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
			return c.Next()
		})
	}

	// Routes
	app.Post("/google/login", handlers.GoogleLogin)
	app.Get("/google/callback", handlers.GoogleCallback)
	app.Get("/logout", handlers.Logout)

	app.Static("/", "./client")

	log.Fatal(app.Listen(os.Getenv("HOST") + ":" + os.Getenv("PORT")))
}
