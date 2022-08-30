package repo_corporatesubs

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/go-common-module/common/logger"
	gormpg "github.com/mindtera/go-common-module/common/v2/configuration/gorm"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CorporateSubsRepoImpl struct {
	sugar     logger.CustomLogger
	pgConfig  gormpg.PostgresConfig
	assertSvc assert.Assert
}

func NewCorporateSubscriptionRepository(sugar logger.CustomLogger,
	pgConfig gormpg.PostgresConfig, assertSvc assert.Assert) CorporateSubsRepo {

	param := pgConfig.GetParam()
	if param.Automigrate {
		sugar.WithContext(context.Background()).Info("creating Corporate subscription relation table")
		client := pgConfig.GetClient()
		client.AutoMigrate(&entity.CorporateSubscriptionEntity{})
		client.AutoMigrate(&entity.CorporateSubscriptionHistory{})
	}

	return &CorporateSubsRepoImpl{
		sugar:     sugar,
		pgConfig:  pgConfig,
		assertSvc: assertSvc,
	}
}

func (c *CorporateSubsRepoImpl) GetCorporateSubscriptionByCorporateID(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error) {

	timeNow := time.Now()
	// var abc map[string]any
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateSubscriptionEntity{}).
		Select(`corporate_subscription_entities.id, corporate_subscription_entities.corporate_id, 
			corporate_subscription_entities.subscription_plan_id, corporate_subscription_entities.subscription_code, 
			corporate_subscription_entities.employee_capacity, corporate_subscription_entities.start_period, 
			corporate_subscription_entities.end_period, corporate_subscription_entities.record_flag`).
		Where("corporate_id = ? AND (start_period <= ? AND end_period >= ?)",
			corpSubs.CorporateID, timeNow, timeNow).
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Preload("CorporateDetail", common.ACTIVE_RECORD_FLAG_QUERY,
			func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name")
			}).
		Find(corpSubs)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to get:%v", err.Error())
		return err
	}

	return err
}

func (c *CorporateSubsRepoImpl) GetCorporateSubscriptionByIDAndCorporateID(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error) {

	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Model(&entity.CorporateSubscriptionEntity{}).
		Select(`id, corporate_id, subscription_code,subscription_plan_id,
		 employee_capacity, start_period, end_period, record_flag`).
		Where("corporate_id = ? AND id = ?",
			corpSubs.CorporateID, corpSubs.ID).
		Where(common.ACTIVE_RECORD_FLAG_QUERY).
		Find(corpSubs)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to get:%v", err.Error())
		return err
	}

	return err
}

func (c *CorporateSubsRepoImpl) UpsertCorporateSubscription(ctx context.Context, corpSubs *entity.CorporateSubscriptionEntity) (err error) {
	var subsHistory entity.CorporateSubscriptionHistory
	subsHistory.BuilderFromCorporateSubscription(*corpSubs)

	if !c.assertSvc.IsUUIDEmpty(corpSubs.ID.String()) {
		subsHistory.Type = "UPDATED"
	}

	// check corporate
	corporate := struct {
		ID uuid.UUID `json:"id" gorm:"id"`
	}{}
	tx := c.pgConfig.GenerateTransaction(ctx)
	tx.Table("corporate_entities").
		Select("id").
		Where("id = ?", corpSubs.CorporateID).
		First(&corporate)

	if c.assertSvc.IsUUIDEmpty(corporate.ID.String()) {
		tx.Rollback()
		err = errors.New("corporate is not found")
		return err
	}

	tx.Model(&entity.CorporateSubscriptionEntity{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"}},
			UpdateAll: true}).
		Create(corpSubs)

	// append corporate history subscription
	subsHistory.SubscriptionID = corpSubs.ID
	tx.Model(&entity.CorporateSubscriptionHistory{}).
		Create(&subsHistory)

	if err = tx.Error; err != nil {
		tx.Rollback()
		c.sugar.WithContext(ctx).Errorf("error to get:%v", err.Error())
		return err
	}
	return err
}
