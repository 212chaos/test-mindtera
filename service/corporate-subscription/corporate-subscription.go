package service_corporatesubscription

import (
	"context"

	"github.com/mindtera/corporate-service/entity"
)

type CorporateSubscriptionService interface {
	// get corporate subscription service
	GetCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error)

	// upsert corporate subscription service
	UpsertCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error)

	// scheduler callback service corporate subscription service
	ExpiredCallbackCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error)
}
