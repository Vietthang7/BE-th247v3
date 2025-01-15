package main

import (
	"fmt"
	"intern_247/config"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("ok")
		return c.SendString("Hello, World!") // Trả về một thông báo cho client
	})
	log.Fatal(app.Listen(":3000"))
}
