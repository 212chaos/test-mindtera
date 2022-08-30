package service_corporatevalue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"

	commonmodel "github.com/mindtera/go-common-module/common/v2/model"
)

func (c *CorporateValueServiceImpl) getCorporateAssessment(ctx context.Context, corporateValue *[]model.CorporateValueModel) (err error) {
	// call parser
	chanSMM := make(chan map[string]entity.CorporateSMM)
	chanErr := make(chan error)
	go func() {
		ctx_ := context.WithValue(ctx, commonmodel.TRANSACTION_KEY, nil)
		ctx_ = context.WithValue(ctx_, string(commonmodel.TRANSACTION_KEY), nil)

		hm, er := c.GetCorporateSMMHash(ctx_)
		chanErr <- er
		chanSMM <- hm
	}()

	// initiate key
	key := "CORPORATE_ASSESSMENT"
	tempMap := map[string][]entity.CorporateValueEntity{}
	if er := c.redisService.Get(ctx, key, &tempMap); er != nil {
		//get from db instead
		c.sugar.WithContext(ctx).Info("get corporate assessment from db")
		var corporateAssessment []entity.CorporateAssessmentRelation
		if err := c.corpAssessRepo.GetCorporateAssessmentRelation(ctx, &corporateAssessment); err != nil {
			return err
		}
		// do populating logic
		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for _, v := range corporateAssessment {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, val entity.CorporateAssessmentRelation) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				// generate url
				iconLink, iconOutlineLink := c.GetValueIcons((*val.CorporateValueDetail).Code)
				// populate url link
				(*val.CorporateValueDetail).IconLink = iconLink
				(*val.CorporateValueDetail).IconLinkOutline = iconOutlineLink
				(*val.CorporateValueDetail).TransformNames()
				tempMap[val.ParentSMM] = append(tempMap[val.ParentSMM], (*val.CorporateValueDetail))
			}(&wg, v)
		}
		wg.Wait()
	}

	// call channel
	if <-chanErr != nil {
		c.sugar.WithContext(ctx).Errorf("error from channel:%v", err.Error())
		return errors.New("error in channel callback")
	}

	// set to redis
	go func() {
		if len(tempMap) <= 0 {
			c.sugar.WithContext(ctx).Errorf("corporate assessment is empty")
			return
		}
		if err = c.redisService.Set(context.Background(), key, tempMap, time.Duration(common.COMMON_TTL*int(time.Minute))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when set data to redis:%v", key)
		}
	}()

	// transform to corporate model
	smmHash := <-chanSMM
	wg2, mx2 := c.commonService.GenerateWaitGroupAndMutex()
	var corporateValueTemp []model.CorporateValueModel
	for i, j := range tempMap {
		wg2.Add(1)
		go func(wg_ *sync.WaitGroup, k string, v []entity.CorporateValueEntity) {
			defer c.commonService.EndWaitGroupAndMutex(wg_, mx2)
			mx2.Lock()
			(corporateValueTemp) = append((corporateValueTemp), model.CorporateValueModel{
				ParentSMM:             k,
				SMMDetail:             smmHash[k],
				CorporateValueRecords: v,
			})
		}(&wg2, i, j)
	}
	wg2.Wait()

	*corporateValue = corporateValueTemp
	return err
}
