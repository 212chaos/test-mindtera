package repo_publicusersubscription

import (
	"context"
	"strings"

	"github.com/google/uuid"
	ce "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSubscriptionRepositoryImpl struct {
	sugar    logger.CustomLogger
	pgConfig gormpg.PostgresConfig
}

func NewUserSubscriptionRepository(
	sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig) UserSubscriptionRepo {
	return &UserSubscriptionRepositoryImpl{
		sugar:    sugar,
		pgConfig: pgConfig,
	}
}

// activate public user subscription
func (p *UserSubscriptionRepositoryImpl) ActiveUserSubscription(ctx context.Context, userSubscription *[]ce.Subscription) (err error) {
	// generate user id array
	userArrChan := make(chan []uuid.UUID)
	go func() {
		var arr []uuid.UUID
		for _, v := range *userSubscription {
			arr = append(arr, v.UserID)
		}
		userArrChan <- arr
	}()

	// generate transaction first
	tx := p.pgConfig.GenerateTransaction(ctx)
	tx.Model(ce.Subscription{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"}},
			UpdateAll: true}).
		Create(userSubscription)

	arr := <-userArrChan
	tx.Table("users").
		Where("id IN ?", arr).
		Update("account_type", "PAID")

	if err = tx.Error; err != nil {
		p.sugar.WithContext(ctx).Error(err.Error())

		tx.Rollback()
		return err
	}

	return err
}

// deactivate public user subscription
func (p *UserSubscriptionRepositoryImpl) UpdateUserSubscriptionStatus(ctx context.Context, userSubscription *ce.Subscription, status string) (err error) {
	// generate transaction first
	tx := p.pgConfig.GenerateTransaction(ctx)
	tx.Model(ce.Subscription{}).
		Where("user_id = ? AND subscription_plan_id = ? ", userSubscription.UserID, userSubscription.SubscriptionPlanID).
		Updates(map[string]any{"record_flag": status})

	p.sugar.WithContext(ctx).Infof("deactivate user subscription:%v", userSubscription.UserID)
	if !strings.EqualFold(status, "PAID") {
		tx.Table("users").
			Where("id = ?", userSubscription.UserID).
			Update("account_type", "FREE")
	}

	if err = tx.Error; err != nil {
		p.sugar.WithContext(ctx).Error(err.Error())

		tx.Rollback()
		return err
	}
	return err
}

// update employee subscription status with transaction
func (p *UserSubscriptionRepositoryImpl) UpdateUserSubscriptionStatusWithTx(ctx context.Context, tx *gorm.DB, userSubscription *ce.Subscription, status string) (err error) {
	// generate transaction first
	tx.Model(ce.Subscription{}).
		Where("user_id = ? AND subscription_plan_id = ? ", userSubscription.UserID, userSubscription.SubscriptionPlanID).
		Update("record_flag", status)

	p.sugar.WithContext(ctx).Infof("deactivate user subscription:%v", userSubscription.UserID)
	if !strings.EqualFold(status, "PAID") {
		tx.Table("users").
			Where("id = ?", userSubscription.UserID).
			Update("account_type", "FREE")
	}
	return err
}
