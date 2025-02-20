package repo

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"intern_247/consts"
	"intern_247/models"
)

type TokenData struct {
	ID          uuid.UUID
	RoleId      int64
	FullName    string
	Username    string
	Position    int64
	StudentType int64
	CenterId    uuid.UUID
	BranchId    *uuid.UUID
	Raw         string
}

func GetTokenData(c *fiber.Ctx) (token TokenData, err error) {
	var (
		ok      bool
		user    models.User
		student Student
		roleId  float64
	)
	if roleId, ok = c.Locals("role_id").(float64); !ok {
		return token, errors.New("error get role_id from context")
	} else {
		token.RoleId = int64(roleId)
		if token.RoleId == consts.Student {
			if student, ok = c.Locals("student").(Student); !ok {
				return token, errors.New("error get student from context")
			} else {
				token.ID = student.ID
				token.FullName = student.FullName
				token.Username = student.Username
				token.CenterId = student.CenterId
				token.BranchId = student.BranchId
				token.StudentType = student.Type
			}
		} else if token.RoleId == consts.Root {

		} else {
			if user, ok = c.Locals("user").(models.User); !ok {
				return token, errors.New("error get user from context")
			} else {
				token.ID = user.ID
				token.FullName = user.FullName
				token.Username = user.Username
				token.CenterId = *user.CenterId
				token.BranchId = user.BranchId
				token.Position = user.Position
			}
		}
	}
	token.Raw = c.Locals("user_token").(*jwt.Token).Raw
	return token, nil
}
