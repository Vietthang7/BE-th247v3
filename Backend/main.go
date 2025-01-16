package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	app2 "intern_247/app"
	"log"
)

func main() {
	app2.Setup()
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("ok")
		return c.SendString("Hello, World!") // Trả về một thông báo cho client
	})
	log.Fatal(app.Listen(":3000"))
}
