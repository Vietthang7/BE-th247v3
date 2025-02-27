package consts

import "github.com/google/uuid"

type DeleteInput struct {
	Id      uuid.UUID `json:"id"`
	ClassId uuid.UUID `json:"classId"`
}

type DeleteMultipleInput struct {
	Ids []uuid.UUID `json:"ids"`
}
