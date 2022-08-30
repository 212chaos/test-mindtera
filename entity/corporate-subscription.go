package entity

import (
	"time"

	"github.com/google/uuid"
	ce "github.com/mindtera/dashboard-auth/entity"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateSubscriptionEntity struct {
	ID                 uuid.UUID `json:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()"`
	CorporateID        uuid.UUID `json:"corporate_id" gorm:"type:uuid;index" binding:"required"`
	SubscriptionCode   string    `json:"subscription_code" binding:"required"`
	SubscriptionPlanID uuid.UUID `json:"subscription_plan_id" gorm:"type:uuid" binding:"required"`
	EmployeeCapacity   int       `json:"employee_capacity" binding:"required"`
	StartPeriod        time.Time `json:"start_period" gorm:"type:timestamp" binding:"required"`
	EndPeriod          time.Time `json:"end_period" gorm:"type:timestamp" binding:"required"`
	c.DefaultColumn
	RecordFlag string `json:"record_flag" gorm:"column:record_flag;index;default:'ACTIVE'"`

	CorporateDetail *ce.CorporateEntity `json:"corporate_detail,omitempty" gorm:"references:CorporateID;foreignKey:ID"`
}
