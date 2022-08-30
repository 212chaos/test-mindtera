package repo_publicusersubscription

import (
	"context"

	ce "github.com/mindtera/consumer-service/model"
	"gorm.io/gorm"
)

type UserSubscriptionRepo interface {
	ActiveUserSubscription(ctx context.Context, userSubscription *[]ce.Subscription) (err error)
	UpdateUserSubscriptionStatus(ctx context.Context, userSubscription *ce.Subscription, status string) (err error)

	// update employee subscription status with transaction
	UpdateUserSubscriptionStatusWithTx(ctx context.Context, tx *gorm.DB, userSubscription *ce.Subscription, status string) (err error)
}
