package models

import "github.com/google/uuid"

type ReqIds struct {
	Ids []*uuid.UUID `json:"ids"`
}
type ReqInputIds struct {
	Ids []uuid.UUID `json:"ids"`
}
