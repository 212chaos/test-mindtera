package service_usersubscription

import (
	"context"
	"sync"

	ce "github.com/mindtera/consumer-service/model"
	"github.com/mindtera/corporate-service/entity"
	usersubs "github.com/mindtera/corporate-service/repository/user-subscription"
	"github.com/mindtera/go-common-module/common/logger"
	assert "github.com/mindtera/go-common-module/common/v2/service/assert"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
)

type UserSubsServiceImpl struct {
	sugar        logger.CustomLogger
	userSubsRepo usersubs.UserSubscriptionRepo
	commonSvc    commonsvc.CommonService
	assert       assert.Assert
}

func NewPublicUserSubscriptionService(sugar logger.CustomLogger,
	userSubsRepo usersubs.UserSubscriptionRepo,
	assert assert.Assert,
	commonSvc commonsvc.CommonService) UserSubsService {
	return &UserSubsServiceImpl{
		sugar:        sugar,
		userSubsRepo: userSubsRepo,
		commonSvc:    commonSvc,
		assert:       assert,
	}
}

// activate and create new public user subscription
func (p *UserSubsServiceImpl) ActiveUserSubsService(ctx context.Context, corporateEmployee *[]entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity) (err error) {
	// user subscription payload
	var userSubscriptions []ce.Subscription

	// loops over payload
	wg, mx := p.commonSvc.GenerateWaitGroupAndMutex()
	for _, val := range *corporateEmployee {
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, v entity.CorporateEmployeeEntity) {
			defer p.commonSvc.EndWaitGroupAndMutex(wg_, mx)
			// lock the payload
			mx.Lock()
			if !p.assert.IsUUIDEmpty(v.PublicUserID.String()) {
				// generate payload for the table first
				userSubscription := ce.Subscription{
					UserID:             v.PublicUserID,
					SubscriptionPlanID: &subscriptionDetail.SubscriptionPlanID,
					StartDate:          &subscriptionDetail.StartPeriod,
					EndDate:            &subscriptionDetail.EndPeriod,
					DefaultColumns:     ce.DefaultColumns{RecordFlag: "ACTIVE"}}
				userSubscriptions = append(userSubscriptions, userSubscription)
			}
		}(&wg, val)
	}
	wg.Wait()

	// upsert the payload created
	p.sugar.WithContext(ctx).Infof("activating public user subscription for corporate: %v___%v", subscriptionDetail.CorporateID, subscriptionDetail.ID)
	if len(userSubscriptions) > 0 {
		if err = p.userSubsRepo.ActiveUserSubscription(ctx, &userSubscriptions); err != nil {
			p.sugar.WithContext(ctx).Errorf("error when upserting public user subscription:%v", err.Error())
			return err
		}
	}

	return err
}

// update the existing public user subscription to DEACTIVATE or ACTIVATE
func (p *UserSubsServiceImpl) UpdateUserSubsStatusService(ctx context.Context, corporateEmployee *entity.CorporateEmployeeEntity, status string, subscriptionDetail *entity.CorporateSubscriptionEntity) (err error) {

	// generate payload for the table first
	userSubscription := ce.Subscription{
		UserID:             corporateEmployee.PublicUserID,
		SubscriptionPlanID: &subscriptionDetail.SubscriptionPlanID,
		StartDate:          &subscriptionDetail.StartPeriod,
		EndDate:            &subscriptionDetail.EndPeriod}

	// updating repository
	p.sugar.WithContext(ctx).Infof("updating employee corporate id:%v status:%v", corporateEmployee.PublicUserID, status)
	if err = p.userSubsRepo.UpdateUserSubscriptionStatus(ctx, &userSubscription, status); err != nil {
		p.sugar.WithContext(ctx).Errorf("error when updating public user subscription:%v", err.Error())
		return err
	}
	return err
}

func (p *UserSubsServiceImpl) DeactivateConcurrentUserFromCorpEmpl(ctx context.Context, emplChan <-chan entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity, status string) {
	go func() {
		for i := range emplChan {
			if err := p.UpdateUserSubsStatusService(ctx, &i, status, subscriptionDetail); err != nil {
				p.sugar.WithContext(ctx).Errorf("error when updating corp empl with corp id:%v subs id:%v empl id:%v",
					subscriptionDetail.CorporateID,
					subscriptionDetail.ID,
					i.PublicUserID)
			}
		}
	}()
}

func (p *UserSubsServiceImpl) ActivateConcurrentUserFromCorpEmpl(ctx context.Context, emplChan <-chan entity.CorporateEmployeeEntity, subscriptionDetail *entity.CorporateSubscriptionEntity) {
	go func() {
		var empls []entity.CorporateEmployeeEntity
		for i := range emplChan {
			empls = append(empls, i)
		}
		if err := p.ActiveUserSubsService(ctx, &empls, subscriptionDetail); err != nil {
			p.sugar.WithContext(ctx).Errorf("error when updating corp empl with corp id:%v subs id:%v",
				subscriptionDetail.CorporateID,
				subscriptionDetail.ID)
		}
	}()
}
