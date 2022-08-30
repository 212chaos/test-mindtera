package entity

import (
	"github.com/google/uuid"
	ca "github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateValueRelation struct {
	ID               uuid.UUID `gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	CorporateID      uuid.UUID `json:"corporate_id" gorm:"type:uuid;uniqueIndex:corporate_relation_unique_2"`
	CorporateValueID uuid.UUID `json:"corporate_value_id" gorm:"type:uuid;uniqueIndex:corporate_relation_unique_2"`
	c.DefaultColumn
	RecordFlag string `json:"record_flag" gorm:"column:record_flag;index;default:'ACTIVE'"`

	// CorporateDetail ca.CorporateEntity
	CorporateValueDetail *CorporateValueEntity `json:"corporate_value_detail,omitempty" gorm:"references:CorporateValueID;foreignKey:ID"`
	CorporateDetail      *ca.CorporateEntity   `json:"corporate_detail,omitempty" gorm:"references:CorporateID;foreignKey:ID"`
}
