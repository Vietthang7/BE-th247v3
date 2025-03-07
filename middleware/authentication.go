package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"intern_247/app"
	"intern_247/consts"
	"intern_247/models"
	"intern_247/repo"
	"intern_247/utils"
	"strings"
)

func AdminAuthentication(c *fiber.Ctx) error {
	path := c.Path()
	if strings.Contains(path, "api/admin/login") {
		return c.Next()
	} else {
		return jwtware.New(jwtware.Config{
			SigningKey:   jwtware.SigningKey{Key: []byte(app.Config("SECRET_KEY"))},
			ErrorHandler: handleErrorResponse,
			SuccessHandler: func(c *fiber.Ctx) error {
				return handleCheckAdminClaims(c)
			},
			ContextKey: "user_token"})(c)
	}
}
func handleErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT", "error": consts.ERROR_UNAUTHORIZED})
}

func handleCheckAdminClaims(c *fiber.Ctx) error {
	u := c.Locals("user_token").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	userId, ok := claims["user_id"].(string)
	if !ok {
		logrus.Debug("claims user_id failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT", "error": consts.ERROR_UNAUTHORIZED})
	}
	roleId, ok := claims["role_id"].(float64)
	if !ok {
		logrus.Debug("claims role_id failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT", "error": consts.ERROR_UNAUTHORIZED})
	}
	if userId == consts.USER_ROOT && roleId == consts.Root {
		c.Locals("user", models.User{
			RoleId:   consts.Root,
			FullName: consts.USER_ROOT,
		})
		c.Locals("role_id", roleId)
		return c.Next()
	}
	user_id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT", "error": consts.ERROR_UNAUTHORIZED})
	}
	if roleId == consts.Student {
		var student repo.Student
		if err = student.First("id = ?", []interface{}{user_id}); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Can not find student", "error": consts.GetFailed})
		}
		if !utils.IsVerifiedEmail(&student.EmailVerified) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Can not verify email", "error": consts.EmailIsNotVerified})
		}
		c.Locals("student", student)
		c.Locals("role_id", roleId)
		return c.Next()
	} else {
		user, row, err := repo.GetUserByID(user_id)
		//jsonData, _ := json.MarshalIndent(user, "", "  ")
		//fmt.Println(string(jsonData))
		if err != nil || row < 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Can not find user", "error": consts.GetFailed})
		} else {
			//user not active
			if !utils.IsActiveData(user.IsActive) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Can not find user", "error": consts.UserIsInactive})
			}
			//user not verify
			if !utils.IsVerifiedEmail(&user.EmailVerified) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Can not verify email", "error": consts.EmailIsNotVerified})
			}
			c.Locals("user", user)
			c.Locals("role_id", roleId)
			return c.Next()
		}
	}
}
func Gate(Subject, Action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Subject == "" || Action == "" {
			return c.Next()
		}
		if roleId, ok := c.Locals("role_id").(float64); ok {
			if roleId == consts.Student {
				return c.Next()
			} else {
				if user, ok := c.Locals("user").(models.User); ok {
					if user.RoleId == consts.Root || user.RoleId == consts.CenterOwner {
						return c.Next()
					}
					if repo.HasPermission(user, Subject, Action) {
						return c.Next()
					}
				}
			}
		}
		return c.SendStatus(fiber.StatusForbidden)
	}
}

// trả về một fiber.Handler (middleware), nghĩa là có thể được gắn vào các route trong Fiber.
func Gate2(Action string, Subject ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(Subject) == 0 || Action == "" {
			return c.Next()
		}
		if roleId, ok := c.Locals("role_id").(float64); ok {
			if roleId == consts.Student {
				return c.Next()
			} else {
				if user, ok := c.Locals("user").(models.User); ok {
					if user.RoleId == consts.Root || user.RoleId == consts.CenterOwner {
						return c.Next()
					}
					if repo.HasPermission2(user, Action, Subject...) {
						return c.Next()
					}
				}
			}
		}
		return c.SendStatus(fiber.StatusForbidden)
	}
}
