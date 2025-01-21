package main

import (
	app2 "intern_247/app"
	"intern_247/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app2.Setup()
	app := fiber.New()
	routes.AdminRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
