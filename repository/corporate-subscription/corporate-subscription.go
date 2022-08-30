package repo_corporatesubs

import (
	"context"

	"github.com/mindtera/corporate-service/entity"
)

type CorporateSubsRepo interface {
	GetCorporateSubscriptionByCorporateID(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error)
	UpsertCorporateSubscription(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error)
	GetCorporateSubscriptionByIDAndCorporateID(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error)
}
