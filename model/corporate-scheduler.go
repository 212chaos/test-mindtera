package model

import (
	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/entity"
)

type CorporateSubscriptionScheduler struct {
	CorporateID     uuid.UUID `json:"corporate_id" binding:"required"`
	SubscriptionID  uuid.UUID `json:"subscription_id" binding:"required"`
	CorporatePeriod string    `json:"corporate_period" binding:"required"`
}

func (c *CorporateSubscriptionScheduler) TransformToSubscriptionEntity() entity.CorporateSubscriptionEntity {
	corporateSubs := entity.CorporateSubscriptionEntity{
		ID:          c.SubscriptionID,
		CorporateID: c.CorporateID,
	}

	return corporateSubs
}

func (c *CorporateSubscriptionScheduler) TransformToEmployeeTask() entity.EmployeeTaskEntity {
	taskEntity := entity.EmployeeTaskEntity{
		SubscriptionID: c.SubscriptionID,
		TaskType:       entity.EMPLOYEE_TASK_TYPE(c.CorporatePeriod),
	}

	return taskEntity
}
