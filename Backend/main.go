package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("ok")
		return c.SendString("Hello, World!") // Trả về một thông báo cho client
	})
	log.Fatal(app.Listen(":3000"))
}
