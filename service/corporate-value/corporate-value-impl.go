package service_corporatevalue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mindtera/corporate-service/common"
	"github.com/mindtera/corporate-service/entity"
	"github.com/mindtera/corporate-service/model"
	corpassessrelrepo "github.com/mindtera/corporate-service/repository/corporate-assessment-relation"
	corpsmmrepo "github.com/mindtera/corporate-service/repository/corporate-smm"
	corpvalrepo "github.com/mindtera/corporate-service/repository/corporate-value"
	"github.com/mindtera/go-common-module/common/logger"
	googlecloud "github.com/mindtera/go-common-module/common/v2/configuration/google-cloud"
	commonsvc "github.com/mindtera/go-common-module/common/v2/service/common"
	redissvc "github.com/mindtera/go-common-module/common/v2/service/redis"
)

type CorporateValueServiceImpl struct {
	sugar          logger.CustomLogger
	gStorage       googlecloud.GoogleStorageConfig
	corpValRepo    corpvalrepo.CorporateValueRepository
	corpAssessRepo corpassessrelrepo.CorporateAssessmentRelationRepository
	corpSMMRepo    corpsmmrepo.CorporateSMMRepository
	redisService   redissvc.RedisSvc
	commonService  commonsvc.CommonService
}

func NewCorporateValueService(sugar logger.CustomLogger,
	gStorage googlecloud.GoogleStorageConfig,
	corpValRepo corpvalrepo.CorporateValueRepository,
	corpAssessRepo corpassessrelrepo.CorporateAssessmentRelationRepository,
	corpSMMRepo corpsmmrepo.CorporateSMMRepository,
	redisService redissvc.RedisSvc,
	commonService commonsvc.CommonService) CorporateValueService {
	return &CorporateValueServiceImpl{
		sugar:          sugar,
		gStorage:       gStorage,
		corpValRepo:    corpValRepo,
		corpAssessRepo: corpAssessRepo,
		corpSMMRepo:    corpSMMRepo,
		redisService:   redisService,
		commonService:  commonService,
	}
}

func (c *CorporateValueServiceImpl) GetCorporateValueService(ctx context.Context, corporateValues *[]entity.CorporateValueEntity) (err error) {
	key := "CORPORATE_VALUES"

	//get redis value
	if err = c.redisService.Get(ctx, key, corporateValues); err != nil {
		//redis is empty
		c.sugar.WithContext(ctx).Infof("redis value is empty: %v", key)

		// fetch from database
		if err = c.corpValRepo.GetCorporateValue(ctx, corporateValues); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when fetching database data:%v", err.Error())
			return err
		}

		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for i, v := range *corporateValues {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, idx int, val entity.CorporateValueEntity) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				iconUrl, iconOutlineUrl := c.GetValueIcons(val.Code)
				// populate icons
				(*corporateValues)[idx].IconLink = iconUrl
				(*corporateValues)[idx].IconLinkOutline = iconOutlineUrl
				// populate localization
				(*corporateValues)[idx].TransformNames()
			}(&wg, i, v)
		}
		wg.Wait()
	}

	go func() {
		// check data len
		if err = c.redisService.Set(context.Background(), key, corporateValues, time.Duration(common.COMMON_TTL*int(time.Minute))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when set cache data:%v", err.Error())
		}
	}()

	return err
}

func (c *CorporateValueServiceImpl) GetCorporateValueBySMMService(ctx context.Context, corporateValues *[]model.CorporateValueModel) (err error) {
	key := "CORPORATE_VALUE_BY_ASSESSMENT"

	// getting value from redis first
	if err := c.redisService.Get(ctx, key, corporateValues); err != nil {
		c.sugar.WithContext(ctx).Info("get corporate assessment from processing db")
		if err := c.getCorporateAssessment(ctx, corporateValues); err != nil {
			return err
		}
	}
	// set to redis
	if corporateValues == nil || len(*corporateValues) <= 0 {
		err = errors.New("corporate values is not found")
		return err
	}
	go func() {
		if err := c.redisService.Set(context.Background(), key, corporateValues,
			time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error when set data to redis:%v", key)
		}
	}()
	return err
}

func (c *CorporateValueServiceImpl) GetCorporateSMMHash(ctx context.Context) (map[string]entity.CorporateSMM, error) {
	key := "CORPORATE_SMM_HASH"
	corporateSMMHash := make(map[string]entity.CorporateSMM)

	if er := c.redisService.Get(ctx, key, &corporateSMMHash); er != nil {
		c.sugar.WithContext(ctx).Info("fetch corporate smm from database")
		var corporateSMM []entity.CorporateSMM
		if err := c.corpSMMRepo.GetCorporateSMM(ctx, &corporateSMM); err != nil {
			return corporateSMMHash, err
		}

		wg, mx := c.commonService.GenerateWaitGroupAndMutex()
		for _, k := range corporateSMM {
			wg.Add(1)
			go func(wg_ *sync.WaitGroup, v entity.CorporateSMM) {
				defer c.commonService.EndWaitGroupAndMutex(wg_, mx)
				mx.Lock()
				//url creation
				v.Icon = fmt.Sprintf("https://storage.googleapis.com/%v/%v",
					common.ICON_BUCKET, v.Icon)
				// populate localization
				v.TransformNames()
				v.TransformDescription()
				corporateSMMHash[v.Code] = v

			}(&wg, k)
		}
		wg.Wait()
	}
	// set redis cache for corporate SMM
	if len(corporateSMMHash) <= 0 {
		return corporateSMMHash, errors.New("smm is not found")
	}
	go func() {
		if err := c.redisService.Set(context.Background(), key, corporateSMMHash,
			time.Duration(common.COMMON_TTL*int(time.Second))); err != nil {
			c.sugar.WithContext(ctx).Errorf("error to set corporate smm cache")
		}
	}()

	return corporateSMMHash, nil
}

// get corporate values detail from id query
func (c *CorporateValueServiceImpl) GetCorporateValueServiceByID(ctx context.Context, cids []uuid.UUID, corporateValues *[]entity.CorporateValueEntity) (err error) {
	// should be cached, but now direct query from database only
	if err = c.corpValRepo.GetCorporateValueByID(ctx, cids, corporateValues); err != nil {
		return err
	}
	return err
}

func (c *CorporateValueServiceImpl) GetValueIcons(englishName string) (iconUrl, iconOutlineUrl string) {
	editedEnglishName := strings.ReplaceAll(strings.ToUpper(englishName), " ", "_") + ".png"
	// generate icon url
	iconUrl = fmt.Sprintf("https://storage.googleapis.com/%v/%v/%v", common.ICON_BUCKET,
		common.ICON_VALUE_FOLDER, editedEnglishName)
	iconOutlineUrl = fmt.Sprintf("https://storage.googleapis.com/%v/%v/%v", common.ICON_BUCKET,
		common.ICON_OUTLINE_VALUE_FOLDER, editedEnglishName)

	return iconUrl, iconOutlineUrl
}
