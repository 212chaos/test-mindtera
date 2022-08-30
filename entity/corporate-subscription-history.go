package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	c "github.com/mindtera/go-common-module/common/entity"
)

type CorporateSubscriptionHistory struct {
	ID                 uuid.UUID `json:"id" gorm:"primaryKey; unique; type:uuid; column:id; not null; default:uuid_generate_v4()"`
	CorporateID        uuid.UUID `json:"corporate_id" gorm:"type:uuid;index" binding:"required"`
	SubscriptionID     uuid.UUID `json:"subscription_id"`
	SubscriptionCode   string    `json:"subscription_code" binding:"required"`
	SubscriptionPlanID uuid.UUID `json:"subscription_plan_id" gorm:"type:uuid" binding:"required"`
	EmployeeCapacity   int       `json:"employee_capacity" binding:"required"`
	StartPeriod        time.Time `json:"start_period" gorm:"type:timestamp" binding:"required"`
	EndPeriod          time.Time `json:"end_period" gorm:"type:timestamp" binding:"required"`
	c.DefaultColumn
	RecordFlag string `json:"record_flag" gorm:"column:record_flag;index;default:'ACTIVE'"`
	Type       string `json:"type"`
}

func (c *CorporateSubscriptionHistory) BuilderFromCorporateSubscription(ce CorporateSubscriptionEntity) {
	b, _ := json.Marshal(ce)
	json.Unmarshal(b, c)

	// append new value
	c.ID = uuid.New()
	c.Type = "CREATED"
}
