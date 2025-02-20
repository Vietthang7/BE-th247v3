package routes

import (
	"github.com/gofiber/fiber/v2"
	"intern_247/controllers"
)

func AuthRoute(auth fiber.Router) {
	auth.Post("/login", controllers.NewLoginGetToken)
	auth.Post("/forgot", controllers.ForgotPwd)
	//verify email
	auth.Post("/verify_email", controllers.VerifyEmailOTP)
	//auth.Post("/verify_token",)
	auth.Post("/register", controllers.Register)
	//resend otp
	auth.Post("/send-otp", controllers.ResendOTP)

}
