package service_corporatesubscription

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	corporatesubsrepo "github.com/mindtera/corporate-service/repository/corporate-subscription"
	schedulerclient "github.com/mindtera/go-common-module/common/client/scheduler-client"
	"github.com/mindtera/go-common-module/common/logger"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
	schedulerentity "github.com/mindtera/scheduler-service/entity"
)

type CorporateSubscriptionServiceImpl struct {
	sugar        logger.CustomLogger
	assert       assert.Assert
	corpSubsRepo corporatesubsrepo.CorporateSubsRepo
	commonSvc    commonsvc.CommonService
	redisSvc     redissvc.RedisSvc
	scheduler    schedulerclient.SchedulerClient
}

var (
	key = "CORPORATE_SUBSCRIPTION_"
)

func NewCorporateSubscriptionService(sugar logger.CustomLogger,
	corpSubsRepo corporatesubsrepo.CorporateSubsRepo,
	assert assert.Assert,
	commonSvc commonsvc.CommonService,
	scheduler schedulerclient.SchedulerClient,
	redisSvc redissvc.RedisSvc) CorporateSubscriptionService {
	return &CorporateSubscriptionServiceImpl{
		sugar:        sugar,
		assert:       assert,
		corpSubsRepo: corpSubsRepo,
		commonSvc:    commonSvc,
		redisSvc:     redisSvc,
		scheduler:    scheduler,
	}
}

// get corporate subscription service
func (c *CorporateSubscriptionServiceImpl) GetCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error) {
	// check corporate id first
	if c.assert.IsUUIDEmpty(corporateSubscription.CorporateID.String()) {
		err = errors.New("corporate id is empty")
		c.sugar.WithContext(ctx).Errorf("corporate id is empty, cannot get subscription")
		return err
	}
	// try to fetch in redis first
	if err = c.redisSvc.Get(ctx, key+corporateSubscription.CorporateID.String(), corporateSubscription); err != nil { // fetching data from database
		if err = c.corpSubsRepo.GetCorporateSubscriptionByCorporateID(ctx, corporateSubscription); err != nil {
			c.sugar.WithContext(ctx).Errorf("error fetching data subscription with corporate id:%v err:%v", corporateSubscription.CorporateID, err)
			return err
		}
		// check subscription code
		if c.assert.IsEmpty(corporateSubscription.SubscriptionCode) {
			err = errors.New("subscription is not found or not valid")
			c.sugar.WithContext(ctx).Errorf("subscription is invalid for corporate id:%v", corporateSubscription.ID)
			return err
		}
	}
	// update the key
	go func() {
		if err = c.redisSvc.Set(context.Background(), key+corporateSubscription.CorporateID.String(), corporateSubscription,
			time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error to set redis key value corporate id:%v", corporateSubscription.ID)
		}
	}()

	return err
}

// upsert corporate subscription service
func (c *CorporateSubscriptionServiceImpl) UpsertCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error) {
	// concurrent delete redis value
	chanErr := make(chan error)
	go func() {
		var err2 error
		if err2 = c.redisSvc.Delete(context.Background(), key+corporateSubscription.ID.String()+"*"); err2 != nil {
			c.sugar.WithContext(ctx).Errorf("error deleting cache for subscription")
		}
		chanErr <- err2
	}()

	// check start and end time
	if corporateSubscription.EndPeriod.Before(corporateSubscription.StartPeriod) {
		err = errors.New("start time and end time is not valid")
		return err
	}

	if corporateSubscription.EndPeriod.Before(time.Now()) {
		corporateSubscription.RecordFlag = "EXPIRED"
	}

	// add created by and updated by
	un := c.commonSvc.GetUserNameFromContext(ctx)
	if c.assert.IsUUIDEmpty(corporateSubscription.ID.String()) {
		corporateSubscription.CreatedBy = un
	} else {
		corporateSubscription.UpdatedBy = un
	}

	if err = c.corpSubsRepo.UpsertCorporateSubscription(ctx, corporateSubscription); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when upserting corporate subscription")
		return err
	}

	if strings.EqualFold(corporateSubscription.RecordFlag, "ACTIVE") {
		// scheduler for post and followup
		// send scheduler payload for POST and FOLLOWUP assessment
		schedulerEntity, err := c.generateSchedulerInitialPayload(*corporateSubscription)
		if err != nil {
			c.sugar.WithContext(ctx).Errorf("error generate scheduler payload:%v", err)
			err = errors.New("error processing scheduler payload")
			return err
		}

		// sending payload to scheduler
		for _, val := range schedulerEntity {
			go func(v schedulerentity.SchedulerPayload) {
				if e := c.scheduler.PostScheduledPayload(ctx, v); e != nil {
					c.sugar.WithContext(ctx).Errorf("error sending payload to scheduler:%v", e)
					return
				}
				c.sugar.WithContext(ctx).Infof("success sending payload to scheduler %v_%v", corporateSubscription.CorporateID, v.Message)
			}(val)
		}
	}

	if err = <-chanErr; err != nil {
		return err
	}

	return err
}

// scheduler callback service corporate subscription service
func (c *CorporateSubscriptionServiceImpl) ExpiredCallbackCorporateSubscriptionService(ctx context.Context, corporateSubscription *entity.CorporateSubscriptionEntity) (err error) {
	go c.redisSvc.Delete(context.Background(), key+corporateSubscription.CorporateID.String())

	// get the repository first
	c.sugar.WithContext(ctx).Infof("getting subscription for corporate:%v with id:%v", corporateSubscription.CorporateID, corporateSubscription.ID)
	if err = c.corpSubsRepo.GetCorporateSubscriptionByIDAndCorporateID(ctx, corporateSubscription); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when fetching corporate subscription id:%v", err)
		err = errors.New("error processing payload")
		return err
	}
	corporateSubscription.RecordFlag = "EXPIRED"

	//copy variable
	var corpSubs entity.CorporateSubscriptionEntity
	c.commonSvc.ObjectMapper(corporateSubscription, &corpSubs)

	// updating the status
	c.sugar.WithContext(ctx).Infof("update subscription status for corporate:%v with id:%v", corporateSubscription.CorporateID, corporateSubscription.ID)
	if err = c.UpsertCorporateSubscriptionService(ctx, &corpSubs); err != nil {
		c.sugar.WithContext(ctx).Errorf("error when upserting corporate subscription id:%v", err)
		err = errors.New("error processing payload")
		return err
	}

	return err
}
