package entity

import (
	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateEntity struct {
	ID   uuid.UUID `form:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
	c.DefaultColumn
}
