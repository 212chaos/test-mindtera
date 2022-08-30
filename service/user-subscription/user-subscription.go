package service_usersubscription

import (
	"context"

	"github.com/mindtera/corporate-service/entity"
)

type UserSubsService interface {
	ActiveUserSubsService(ctx context.Context, corporateEmployee *[]entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity) (err error)
	UpdateUserSubsStatusService(ctx context.Context, corporateEmployee *entity.CorporateEmployeeEntity, status string, subscriptionDetail *entity.CorporateSubscriptionEntity) (err error)

	// update user from empl
	DeactivateConcurrentUserFromCorpEmpl(ctx context.Context, emplChan <-chan entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity, status string)
	ActivateConcurrentUserFromCorpEmpl(ctx context.Context, emplChan <-chan entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity)
}
