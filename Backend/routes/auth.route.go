package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func AuthRoute(auth fiber.Router) {
	auth.Post("/login", controllers.NewLoginGetToken)
	//verify email
	auth.Post("/verify_email", controllers.VerifyEmailOTP)
}
